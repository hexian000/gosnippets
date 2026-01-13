// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package hlistener

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mockListener is a mock implementation of net.Listener for testing
type mockListener struct {
	acceptCh chan net.Conn
	errCh    chan error
	addr     net.Addr
	closed   atomic.Bool
}

func newMockListener() *mockListener {
	return &mockListener{
		acceptCh: make(chan net.Conn, 10),
		errCh:    make(chan error, 1),
		addr:     &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080},
	}
}

func (m *mockListener) Accept() (net.Conn, error) {
	select {
	case conn := <-m.acceptCh:
		return conn, nil
	case err := <-m.errCh:
		return nil, err
	}
}

func (m *mockListener) Close() error {
	m.closed.Store(true)
	close(m.acceptCh)
	return nil
}

func (m *mockListener) Addr() net.Addr {
	return m.addr
}

// mockConn is a mock implementation of net.Conn for testing
type mockConn struct {
	closed atomic.Bool
}

func newMockConn() *mockConn {
	return &mockConn{}
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockConn) Close() error                       { m.closed.Store(true); return nil }
func (m *mockConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (m *mockConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func (m *mockConn) IsClosed() bool {
	return m.closed.Load()
}

func TestWrap(t *testing.T) {
	mock := newMockListener()
	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)
	if hl == nil {
		t.Fatal("Wrap returned nil")
	}

	if hl.Addr() != mock.Addr() {
		t.Errorf("Expected addr %v, got %v", mock.Addr(), hl.Addr())
	}
}

func TestBasicAccept(t *testing.T) {
	mock := newMockListener()
	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)

	// Send a mock connection
	mockConn := newMockConn()
	mock.acceptCh <- mockConn

	conn, err := hl.Accept()
	if err != nil {
		t.Fatalf("Accept failed: %v", err)
	}
	if conn != mockConn {
		t.Error("Expected to receive the same mock connection")
	}

	total, served := hl.Stats()
	if total != 1 {
		t.Errorf("Expected total=1, got %d", total)
	}
	if served != 1 {
		t.Errorf("Expected served=1, got %d", served)
	}
}

func TestMaxSessionsLimit(t *testing.T) {
	mock := newMockListener()
	var numSessions atomic.Uint32

	config := &Config{
		MaxSessions: 5,
		Stats: func() (uint32, uint32) {
			return numSessions.Load(), 0
		},
	}

	hl := Wrap(mock, config)

	// Set sessions to exceed max
	numSessions.Store(6)

	// Send connections
	for i := 0; i < 10; i++ {
		mockConn := newMockConn()
		mock.acceptCh <- mockConn
	}

	// Try to accept - should be rejected and closed
	var acceptedCount int
	done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			conn, err := hl.Accept()
			if err != nil {
				break
			}
			// Connection was accepted but should be closed immediately
			if conn == nil {
				continue
			}
			if !conn.(*mockConn).IsClosed() {
				acceptedCount++
				conn.Close()
			}
		}
		done <- true
	}()

	time.Sleep(100 * time.Millisecond)

	total, served := hl.Stats()
	if total != 10 {
		t.Errorf("Expected total=10, got %d", total)
	}
	if served != 0 {
		t.Errorf("Expected served=0 (all rejected), got %d", served)
	}
}

func TestFullLimit(t *testing.T) {
	mock := newMockListener()
	var numHalfOpen atomic.Uint32

	config := &Config{
		Full: 3,
		Stats: func() (uint32, uint32) {
			return 0, numHalfOpen.Load()
		},
	}

	hl := Wrap(mock, config)

	// Set half-open to exceed full
	numHalfOpen.Store(5)

	// Send connections
	for i := 0; i < 5; i++ {
		mockConn := newMockConn()
		mock.acceptCh <- mockConn
	}

	// All should be rejected
	go func() {
		for i := 0; i < 5; i++ {
			hl.Accept()
		}
	}()

	time.Sleep(100 * time.Millisecond)

	total, served := hl.Stats()
	if total != 5 {
		t.Errorf("Expected total=5, got %d", total)
	}
	if served != 0 {
		t.Errorf("Expected served=0, got %d", served)
	}
}

func TestStartRateLimit(t *testing.T) {
	mock := newMockListener()
	var numHalfOpen atomic.Uint32

	config := &Config{
		Start: 2,
		Rate:  1.0, // 100% rejection rate when over Start
		Stats: func() (uint32, uint32) {
			return 0, numHalfOpen.Load()
		},
	}

	hl := Wrap(mock, config)

	// Set half-open to exceed start but not full
	numHalfOpen.Store(3)

	// Send connections
	for i := 0; i < 10; i++ {
		mockConn := newMockConn()
		mock.acceptCh <- mockConn
	}

	go func() {
		for i := 0; i < 10; i++ {
			hl.Accept()
		}
	}()

	time.Sleep(100 * time.Millisecond)

	total, served := hl.Stats()
	if total != 10 {
		t.Errorf("Expected total=10, got %d", total)
	}
	// With Rate=1.0, all should be rejected
	if served != 0 {
		t.Errorf("Expected served=0, got %d", served)
	}
}

func TestNoLimit(t *testing.T) {
	mock := newMockListener()

	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)

	// Send multiple connections
	const numConns = 10
	for i := 0; i < numConns; i++ {
		mockConn := newMockConn()
		mock.acceptCh <- mockConn
	}

	var wg sync.WaitGroup
	wg.Add(numConns)

	for i := 0; i < numConns; i++ {
		go func() {
			defer wg.Done()
			conn, err := hl.Accept()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			conn.Close()
		}()
	}

	wg.Wait()

	total, served := hl.Stats()
	if total != numConns {
		t.Errorf("Expected total=%d, got %d", numConns, total)
	}
	if served != numConns {
		t.Errorf("Expected served=%d, got %d", numConns, served)
	}
}

func TestClose(t *testing.T) {
	mock := newMockListener()
	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)

	if err := hl.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if !mock.closed.Load() {
		t.Error("Expected underlying listener to be closed")
	}
}

func TestStats(t *testing.T) {
	mock := newMockListener()
	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)

	// Initially should be 0
	total, served := hl.Stats()
	if total != 0 || served != 0 {
		t.Errorf("Expected stats to be (0,0), got (%d,%d)", total, served)
	}

	// Accept a connection
	mockConn := newMockConn()
	mock.acceptCh <- mockConn

	conn, err := hl.Accept()
	if err != nil {
		t.Fatalf("Accept failed: %v", err)
	}
	conn.Close()

	total, served = hl.Stats()
	if total != 1 || served != 1 {
		t.Errorf("Expected stats to be (1,1), got (%d,%d)", total, served)
	}
}

func TestConcurrentAccept(t *testing.T) {
	mock := newMockListener()
	config := &Config{
		Stats: func() (uint32, uint32) {
			return 0, 0
		},
	}

	hl := Wrap(mock, config)

	const numConns = 20

	var wg sync.WaitGroup
	wg.Add(numConns)

	// Start goroutines to accept connections
	for i := 0; i < numConns; i++ {
		go func(idx int) {
			defer wg.Done()
			conn, err := hl.Accept()
			if err != nil {
				t.Errorf("Accept %d failed: %v", idx, err)
				return
			}
			conn.Close()
		}(i)
	}

	// Give goroutines time to start waiting
	time.Sleep(50 * time.Millisecond)

	// Send connections
	for i := 0; i < numConns; i++ {
		mockConn := newMockConn()
		select {
		case mock.acceptCh <- mockConn:
		case <-time.After(time.Second):
			t.Fatalf("Timeout sending connection %d", i)
		}
	}

	// Wait for all accepts to complete with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for concurrent accepts")
	}

	total, served := hl.Stats()
	if total != numConns {
		t.Errorf("Expected total=%d, got %d", numConns, total)
	}
	if served != numConns {
		t.Errorf("Expected served=%d, got %d", numConns, served)
	}
}
