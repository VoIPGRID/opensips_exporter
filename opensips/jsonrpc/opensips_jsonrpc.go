package jsonrpc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KeisukeYamashita/go-jsonrpc"
	"github.com/VoIPGRID/opensips_exporter/opensips"
)

// JSONRPC holds all the information necessary for handling connections to
// the OpenSIPS Management Interface (targeting version >= 3.0).
type JSONRPC struct {
	url string
}

// New creates a new JSONRPC instance. Pass it the running OpenSIPS'
// HTTP JSON RPC endpoint to connect to.
func New(url string) *JSONRPC {
	return &JSONRPC{
		url: url,
	}
}

// GetStatistics calls the JSON-RPC endpoint and returns the
// statistics OpenSIPS sends back. The targets can be "all", "group:" or
// "name" (e.g. "shmem:" or "rcv_requests").
func (o *JSONRPC) GetStatistics(targets ...string) (map[string]opensips.Statistic, error) {
	rpcClient := jsonrpc.NewRPCClient(o.url)

	// request {"jsonrpc":"2.0","method":"get_statistics","params":[["core:","tm:"]],"id":1}
	resp, err := rpcClient.Call("get_statistics", targets)
	if err != nil {
		return nil, fmt.Errorf("error while getting statistics from JSON-RPC endpoint: %w", err)
	}

	statistics, err := parseStatistics(resp.Result.(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func parseStatistics(response map[string]interface{}) (map[string]opensips.Statistic, error) {
	var res = map[string]opensips.Statistic{}
	for key, value := range response {
		asString := fmt.Sprintf("%s = %s", key, value)
		stat, err := parseStatistic(asString)
		if err != nil {
			return res, fmt.Errorf("error while parsing stat: %w", err)
		}
		res[stat.Name] = stat
	}
	return res, nil
}

func parseStatistic(metric string) (opensips.Statistic, error) {
	var name, module, valueString string
	// Check for OpenSIPS >= 2 metric format
	// i.e.shmem:total_size:: 2147483648
	if metric == "" {
		// There's an empty line in the output since OpenSIPS 2.4.5
		// ignore and continue
		return opensips.Statistic{}, nil
	}
	if strings.Contains(metric, "=") {
		// OpenSIPS < 2 metric format
		// i.e. shmem:total_size = 2147483648
		metricSplit := strings.Split(metric, ":")
		module = metricSplit[0]
		name = strings.Split(strings.Join(metricSplit[1:], ":"), " ")[0]
		i := strings.LastIndex(metric, " ")
		valueString = metric[i+1:]
	} else {
		return opensips.Statistic{}, fmt.Errorf("unknown metric format encountered for: %s", metric)
	}

	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return opensips.Statistic{}, err
	}

	return opensips.Statistic{
		Module: module,
		Name:   name,
		Value:  value,
	}, nil
}
