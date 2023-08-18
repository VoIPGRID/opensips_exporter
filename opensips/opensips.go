package opensips

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"github.com/KeisukeYamashita/go-jsonrpc"
)

const firstLineOK = "200 OK\n"

// OpenSIPS holds all the information necessary for handling connections to
// the OpenSIPS Management Interface (targeting version 1.10).
type OpenSIPS struct {
	socket string
	tmpdir string

	count int64
}

// Statistic holds the module, name and value of a statistic
// as returned by OpenSIPS.
type Statistic struct {
	Module, Name string
	Value        float64
}

// New creates a new OpenSIPS instance. Pass it the running OpenSIPS'
// mi_datagram Unix socket string to connect to. The socket should be
// expressed as a full path to the socket, and the current user should have
// permissions to read from and write to this socket, in addition to write
// access to the folder it's located in (for creating the return socket).
func New(socket string) (*OpenSIPS, error) {
	tmpdir, err := ioutil.TempDir(path.Dir(socket), "opensips_exporter")
	if err != nil {
		return nil, err
	}
	return &OpenSIPS{
		socket: socket,
		tmpdir: tmpdir,
	}, nil
}

// GetStatistics calls the get_statistics management function and returns the
// statistics OpenSIPS sends back. The targets can be "all", "group:" or
// "name" (e.g. "shmem:" or "rcv_requests").
func (o *OpenSIPS) GetStatistics(targets ...string) (map[string]Statistic, error) {
	rpcClient := jsonrpc.NewRPCClient("")
	req := rpcClient.NewRPCRequestObject("get_statistics", targets)
	msg, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := o.roundtrip(msg)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(respBytes))
	decoder.UseNumber()
	rpcResponse := jsonrpc.RPCResponse{}
	err = decoder.Decode(&rpcResponse)
	if err != nil {
		return nil, err
	}
	statistics, err := parseStatistics(rpcResponse.Result.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	return statistics, nil
}

func parseStatistics(response map[string]interface{}) (map[string]Statistic, error) {
	var res = map[string]Statistic{}
	for key, value := range response {
		stat, err := parseStatistic(key, fmt.Sprintf("%s", value))
		if err != nil {
			return res, fmt.Errorf("error while parsing stat: %w", err)
		}
		res[stat.Name] = stat
	}
	return res, nil
}

func parseStatistic(key string, valueString string) (Statistic, error) {
	// OpenSIPS < 2 metric format
	// i.e. shmem:total_size = 2147483648
	metricSplit := strings.Split(key, ":")
	module := metricSplit[0]
	name := metricSplit[1]

	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return Statistic{}, err
	}

	return Statistic{
		Module: module,
		Name:   name,
		Value:  value,
	}, nil
}

func (o *OpenSIPS) roundtrip(request []byte) ([]byte, error) {
	raddr, err := net.ResolveUnixAddr("unixgram", o.socket)
	if err != nil {
		return nil, err
	}
	count := atomic.AddInt64(&o.count, 1)
	laddr, err := net.ResolveUnixAddr("unixgram", path.Join(o.tmpdir, fmt.Sprintf("%d.sock", count)))
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUnixgram("unixgram", laddr)
	if err != nil {
		return nil, err
	}
	defer os.Remove(laddr.Name)
	defer c.Close()
	_, err = c.WriteToUnix(request, raddr)
	if err != nil {
		return nil, err
	}
	err = c.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 65535)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], err
}

// Close tears down all resources created for this OpenSIPS instance.
func (o *OpenSIPS) Close() error {
	err := os.Remove(o.tmpdir)
	return err
}
