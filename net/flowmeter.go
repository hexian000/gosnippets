package net

import (
	"net"
	"sync/atomic"
)

type FlowStats struct {
	Read    atomic.Uint64
	Written atomic.Uint64
}

type meteredConn struct {
	net.Conn
	m *FlowStats
}

func (c *meteredConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.m.Read.Add(uint64(n))
	return
}

func (c *meteredConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.m.Written.Add(uint64(n))
	return
}

func FlowMeter(c net.Conn, m *FlowStats) net.Conn {
	// early validation
	m.Read.Add(0)
	m.Written.Add(0)
	return &meteredConn{c, m}
}
