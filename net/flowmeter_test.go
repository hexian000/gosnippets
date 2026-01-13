// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package net

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// mockConn implements net.Conn for testing purposes
type mockConn struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
}

func newMockConn() *mockConn {
	return &mockConn{
		readBuf:  &bytes.Buffer{},
		writeBuf: &bytes.Buffer{},
	}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuf.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestFlowMeter_Read(t *testing.T) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)

	// Prepare data to read
	testData := []byte("Hello, World!")
	mock.readBuf.Write(testData)

	// Read data
	buf := make([]byte, len(testData))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if n != len(testData) {
		t.Errorf("Read() n = %d, want %d", n, len(testData))
	}

	if !bytes.Equal(buf, testData) {
		t.Errorf("Read() data = %q, want %q", buf, testData)
	}

	// Check stats
	if got := stats.Read.Load(); got != uint64(len(testData)) {
		t.Errorf("FlowStats.Read = %d, want %d", got, len(testData))
	}

	if got := stats.Written.Load(); got != 0 {
		t.Errorf("FlowStats.Written = %d, want 0", got)
	}
}

func TestFlowMeter_Write(t *testing.T) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)

	// Write data
	testData := []byte("Hello, World!")
	n, err := conn.Write(testData)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if n != len(testData) {
		t.Errorf("Write() n = %d, want %d", n, len(testData))
	}

	// Check if data was written to underlying connection
	if !bytes.Equal(mock.writeBuf.Bytes(), testData) {
		t.Errorf("Written data = %q, want %q", mock.writeBuf.Bytes(), testData)
	}

	// Check stats
	if got := stats.Written.Load(); got != uint64(len(testData)) {
		t.Errorf("FlowStats.Written = %d, want %d", got, len(testData))
	}

	if got := stats.Read.Load(); got != 0 {
		t.Errorf("FlowStats.Read = %d, want 0", got)
	}
}

func TestFlowMeter_MultipleOperations(t *testing.T) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)

	// Multiple writes
	data1 := []byte("First write")
	data2 := []byte("Second write")
	
	n1, err := conn.Write(data1)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	n2, err := conn.Write(data2)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	expectedWritten := uint64(n1 + n2)
	if got := stats.Written.Load(); got != expectedWritten {
		t.Errorf("FlowStats.Written = %d, want %d", got, expectedWritten)
	}

	// Prepare data for reading
	readData1 := []byte("First read")
	readData2 := []byte("Second read")
	mock.readBuf.Write(readData1)
	mock.readBuf.Write(readData2)

	// Multiple reads
	buf1 := make([]byte, len(readData1))
	r1, err := conn.Read(buf1)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	buf2 := make([]byte, len(readData2))
	r2, err := conn.Read(buf2)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	expectedRead := uint64(r1 + r2)
	if got := stats.Read.Load(); got != expectedRead {
		t.Errorf("FlowStats.Read = %d, want %d", got, expectedRead)
	}
}

func TestFlowMeter_EOF(t *testing.T) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)

	// Try to read from empty buffer (should get EOF)
	buf := make([]byte, 10)
	n, err := conn.Read(buf)
	if err != io.EOF {
		t.Errorf("Read() error = %v, want EOF", err)
	}

	// Stats should still be updated even on EOF
	if got := stats.Read.Load(); got != uint64(n) {
		t.Errorf("FlowStats.Read = %d, want %d", got, n)
	}
}

func TestFlowMeter_Concurrent(t *testing.T) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)

	// Prepare data for concurrent reads
	for i := 0; i < 100; i++ {
		mock.readBuf.Write([]byte("test data "))
	}

	var wg sync.WaitGroup
	goroutines := 10
	opsPerGoroutine := 10

	// Concurrent writes
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				conn.Write([]byte("test"))
			}
		}()
	}

	// Concurrent reads
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			buf := make([]byte, 10)
			for j := 0; j < opsPerGoroutine; j++ {
				conn.Read(buf)
			}
		}()
	}

	wg.Wait()

	// Verify stats are accurate
	expectedWritten := uint64(goroutines * opsPerGoroutine * 4) // "test" = 4 bytes
	if got := stats.Written.Load(); got != expectedWritten {
		t.Errorf("FlowStats.Written = %d, want %d", got, expectedWritten)
	}

	expectedRead := uint64(goroutines * opsPerGoroutine * 10) // 10 bytes per read
	if got := stats.Read.Load(); got != expectedRead {
		t.Errorf("FlowStats.Read = %d, want %d", got, expectedRead)
	}
}

func TestFlowMeter_NilStats(t *testing.T) {
	mock := newMockConn()

	// Should panic when FlowStats is nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("FlowMeter() should panic with nil FlowStats")
		}
	}()

	FlowMeter(mock, nil)
}

func TestFlowStats_InitialValues(t *testing.T) {
	stats := &FlowStats{}
	
	if got := stats.Read.Load(); got != 0 {
		t.Errorf("FlowStats.Read initial value = %d, want 0", got)
	}

	if got := stats.Written.Load(); got != 0 {
		t.Errorf("FlowStats.Written initial value = %d, want 0", got)
	}
}

func BenchmarkFlowMeter_Read(b *testing.B) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)
	
	data := bytes.Repeat([]byte("x"), 1024)
	buf := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mock.readBuf.Write(data)
		conn.Read(buf)
	}
}

func BenchmarkFlowMeter_Write(b *testing.B) {
	mock := newMockConn()
	stats := &FlowStats{}
	conn := FlowMeter(mock, stats)
	
	data := bytes.Repeat([]byte("x"), 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.Write(data)
	}
}
