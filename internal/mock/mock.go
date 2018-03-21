package mock

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"

	"golang.org/x/sync/errgroup"
)

// Mock is a fake OpenSIPS mi_datagram socket
type Mock struct {
	response []byte
	sleep    time.Duration

	dir  string
	addr *net.UnixAddr
	l    *net.UnixConn
	g    errgroup.Group
}

// New creates a new Mock with given response. Before responding the Mock will
// sleep for a moment.
func New(response []byte, sleep time.Duration) (m *Mock, err error) {
	m = new(Mock)
	m.response = response
	m.sleep = sleep
	m.dir, err = ioutil.TempDir(os.TempDir(), "mock-opensips-")
	if err != nil {
		return
	}
	m.addr, err = net.ResolveUnixAddr("unixgram", path.Join(m.dir, "mock.sock"))
	if err != nil {
		return
	}
	m.l, err = net.ListenUnixgram("unixgram", m.addr)
	if err != nil {
		return
	}
	return
}

// Socket returns the Mock's socket address
func (m *Mock) Socket() string {
	return m.addr.Name
}

// Run handles a given number of requests within the deadline.
func (m *Mock) Run(count int, deadline time.Time) error {
	err := m.l.SetReadDeadline(deadline)
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		buf := make([]byte, 65535)
		_, raddr, err := m.l.ReadFromUnix(buf)
		if err != nil {
			return err
		}
		if raddr == nil {
			return fmt.Errorf("mock.run: got nil raddr")
		}
		m.g.Go(func() error {
			time.Sleep(m.sleep)
			// OpenSIPS responds through an anonymous socket too
			c, err := net.DialUnix("unixgram", nil, raddr)
			if err != nil {
				return err
			}
			_, err = c.Write(m.response)
			if err != nil {
				return err
			}
			return c.Close()
		})
	}
	return m.g.Wait()
}

// Close removes the resources created for Mock.
func (m *Mock) Close() error {
	err := m.l.Close()
	if err != nil {
		return err
	}
	err = os.Remove(m.addr.Name)
	if err != nil {
		return err
	}
	return os.Remove(m.dir)
}
