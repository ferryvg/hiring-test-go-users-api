package http

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type Listener struct {
	ln net.Listener

	maxWaitTime time.Duration
	done        chan struct{}
	connsCount  uint64
	shutdown    uint64
}

// Creates network listener. This listener support graceful Shutdown.
func NewListener(addr string, maxWaitTime time.Duration) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Listener{
		ln:          ln,
		maxWaitTime: maxWaitTime,
		done:        make(chan struct{}),
	}, nil
}

func (ln *Listener) Accept() (net.Conn, error) {
	conn, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}

	atomic.AddUint64(&ln.connsCount, 1)

	return &gracefulConn{Conn: conn, ln: ln}, nil
}

func (ln *Listener) Addr() net.Addr {
	return ln.ln.Addr()
}

func (ln *Listener) Close() error {
	err := ln.ln.Close()
	if err != nil {
		return err
	}

	return ln.waitForClose()
}

func (ln *Listener) waitForClose() error {
	atomic.AddUint64(&ln.shutdown, 1)

	if atomic.LoadUint64(&ln.connsCount) == 0 {
		close(ln.done)
		return nil
	}

	select {
	case <-ln.done:
		return nil
	case <-time.After(ln.maxWaitTime):
		return fmt.Errorf("cannot complete graceful shutdown in %s", ln.maxWaitTime)
	}
}

func (ln *Listener) closeConn() {
	connsCount := atomic.AddUint64(&ln.connsCount, ^uint64(0))
	if atomic.LoadUint64(&ln.shutdown) != 0 && connsCount == 0 {
		close(ln.done)
	}
}

type gracefulConn struct {
	net.Conn
	ln *Listener
}

func (c *gracefulConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.ln.closeConn()

	return nil
}
