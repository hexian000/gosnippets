// The hardened listener
package hlistener

import (
	"math/rand"
	"net"
	"sync/atomic"

	"github.com/hexian000/gosnippets/formats"
	"github.com/hexian000/gosnippets/slog"
)

type Config struct {
	Start, Full uint32
	Rate        int // 0-100
	MaxSessions uint32
	Stats       func() (numSessions uint32, numHalfOpen uint32)
}

type Listener interface {
	net.Listener
	Stats() (total uint64, served uint64)
}

type listener struct {
	l     net.Listener
	c     Config
	stats struct {
		Total  atomic.Uint64
		Served atomic.Uint64
	}
}

func (l *listener) isLimited() bool {
	numSessions, numHalfOpen := l.c.Stats()
	if l.c.MaxSessions > 0 && numSessions >= l.c.MaxSessions {
		return true
	}
	if numHalfOpen >= l.c.Full {
		return true
	}
	if numHalfOpen >= l.c.Start {
		return rand.Intn(100) < l.c.Rate
	}
	return false
}

func (l *listener) Accept() (net.Conn, error) {
	for {
		conn, err := l.l.Accept()
		if err != nil {
			return conn, err
		}
		l.stats.Total.Add(1)
		if l.isLimited() {
			if err := conn.Close(); err != nil {
				slog.Warningf("close: %s", formats.Error(err))
			}
			continue
		}
		l.stats.Served.Add(1)
		return conn, err
	}
}

func (l *listener) Close() error {
	return l.l.Close()
}

func (l *listener) Addr() net.Addr {
	return l.l.Addr()
}

func (l *listener) Stats() (accepted uint64, served uint64) {
	return l.stats.Total.Load(), l.stats.Served.Load()
}

// Wrap the raw listener
func Wrap(l net.Listener, c *Config) Listener {
	return &listener{l: l, c: *c}
}
