// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package net

import (
	"net"
	"sync/atomic"
)

// FlowStats keeps track of the number of bytes read and written.
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

// FlowMeter wraps the given net.Conn and tracks the number of bytes read and written.
func FlowMeter(c net.Conn, m *FlowStats) net.Conn {
	// validate the pointer to FlowStats early
	_ = m.Read.Add(0)
	_ = m.Written.Add(0)
	return &meteredConn{c, m}
}
