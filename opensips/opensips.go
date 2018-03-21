package opensips

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync/atomic"
	"time"
)

const firstLineOK = "200 OK\n"

// OpenSIPS holds all the information necessary for handling connections to
// the OpenSIPS Management Interface (targeting version 1.10).
type OpenSIPS struct {
	socket string
	tmpdir string

	count int64
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
func (o *OpenSIPS) GetStatistics(targets ...string) ([]string, error) {
	msg := []byte(":get_statistics:\n")
	for _, target := range targets {
		msg = append(msg, []byte(target)...)
		msg = append(msg, '\n')
	}
	resp, err := o.roundtrip(msg)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(resp)
	line, err := buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if line != firstLineOK {
		return nil, fmt.Errorf("expected %q, got %q", firstLineOK, line)
	}
	var rv []string
	for err == nil {
		rv = append(rv, line)
		line, err = buf.ReadString('\n')
	}
	return rv[1:], nil
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
