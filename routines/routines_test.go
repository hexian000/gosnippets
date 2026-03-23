// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package routines

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// isPanicError reports whether err is an ErrPanic (value or pointer receiver).
func isPanicError(err error) bool {
	var ep ErrPanic
	var epPtr *ErrPanic
	return errors.As(err, &ep) || errors.As(err, &epPtr)
}

// --- ErrPanic ---

func TestErrPanic(t *testing.T) {
	v := "test panic value"
	ep := ErrPanic{v: v}
	if ep.Panic() != v {
		t.Errorf("ErrPanic.Panic() = %v, want %v", ep.Panic(), v)
	}
	want := "panic: test panic value"
	if ep.Error() != want {
		t.Errorf("ErrPanic.Error() = %q, want %q", ep.Error(), want)
	}
}

// --- Group ---

func TestGroup_Go(t *testing.T) {
	g := NewGroup()
	defer g.Close()
	var mu sync.Mutex
	sum := 0
	for i := 1; i <= 10; i++ {
		i := i
		if err := g.Go(func() {
			mu.Lock()
			sum += i
			mu.Unlock()
		}); err != nil {
			t.Fatalf("Go() error = %v", err)
		}
	}
	if err := g.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
	if sum != 55 {
		t.Errorf("sum = %d, want 55", sum)
	}
}

func TestGroup_Close(t *testing.T) {
	g := NewGroup()
	g.Close()
	if err := g.Go(func() {}); !errors.Is(err, ErrClosed) {
		t.Errorf("Go() after Close() = %v, want ErrClosed", err)
	}
}

func TestGroup_CloseC(t *testing.T) {
	g := NewGroup()
	ch := g.CloseC()
	select {
	case <-ch:
		t.Fatal("CloseC() closed prematurely")
	default:
	}
	g.Close()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("CloseC() not closed after Close()")
	}
}

func TestGroup_Panic(t *testing.T) {
	g := NewGroup()
	defer g.Close()
	if err := g.Go(func() { panic("test panic") }); err != nil {
		t.Fatalf("Go() error = %v", err)
	}
	err := g.Wait()
	if !isPanicError(err) {
		t.Errorf("Wait() = %v, want ErrPanic", err)
	}
}

// --- LimitedGroup ---

func TestLimitedGroup_Go(t *testing.T) {
	const limit = 4
	g := NewLimitedGroup(limit)
	defer g.Close()
	var mu sync.Mutex
	sum := 0
	for i := 1; i <= limit; i++ {
		i := i
		if err := g.Go(func() {
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			sum += i
			mu.Unlock()
		}); err != nil {
			t.Fatalf("Go() error = %v", err)
		}
	}
	if err := g.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
	if sum != 10 {
		t.Errorf("sum = %d, want 10", sum)
	}
}

func TestLimitedGroup_ConcurrencyLimit(t *testing.T) {
	const limit = 2
	g := NewLimitedGroup(limit)
	defer g.Close()
	ready := make(chan struct{})
	block := make(chan struct{})
	for i := 0; i < limit; i++ {
		if err := g.Go(func() {
			ready <- struct{}{}
			<-block
		}); err != nil {
			t.Fatalf("Go() error = %v", err)
		}
	}
	for i := 0; i < limit; i++ {
		<-ready
	}
	if err := g.Go(func() {}); !errors.Is(err, ErrConcurrencyLimit) {
		t.Errorf("Go() beyond limit = %v, want ErrConcurrencyLimit", err)
	}
	close(block)
	if err := g.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
}

func TestLimitedGroup_Close(t *testing.T) {
	g := NewLimitedGroup(4)
	g.Close()
	if err := g.Go(func() {}); !errors.Is(err, ErrClosed) {
		t.Errorf("Go() after Close() = %v, want ErrClosed", err)
	}
}

func TestLimitedGroup_Panic(t *testing.T) {
	g := NewLimitedGroup(2)
	defer g.Close()
	if err := g.Go(func() { panic("limited panic") }); err != nil {
		t.Fatalf("Go() error = %v", err)
	}
	err := g.Wait()
	if !isPanicError(err) {
		t.Errorf("Wait() = %v, want ErrPanic", err)
	}
}

// --- Queue ---

func TestQueue_PushPop(t *testing.T) {
	q := NewQueue[int]()
	for i := 0; i < 5; i++ {
		if !q.Push(i) {
			t.Fatalf("Push(%d) = false, want true", i)
		}
	}
	for i := 0; i < 5; i++ {
		v, ok := q.Pop()
		if !ok {
			t.Fatalf("Pop() ok = false, want true at index %d", i)
		}
		if v != i {
			t.Errorf("Pop() = %d, want %d", v, i)
		}
	}
}

func TestQueue_Len(t *testing.T) {
	q := NewQueue[string]()
	if l := q.Len(); l != 0 {
		t.Errorf("Len() = %d, want 0", l)
	}
	q.Push("a")
	q.Push("b")
	if l := q.Len(); l != 2 {
		t.Errorf("Len() = %d, want 2", l)
	}
	q.Pop()
	if l := q.Len(); l != 1 {
		t.Errorf("Len() = %d, want 1", l)
	}
}

func TestQueue_TryPop_Empty(t *testing.T) {
	q := NewQueue[int]()
	if _, ok := q.TryPop(); ok {
		t.Error("TryPop() on empty queue = true, want false")
	}
}

func TestQueue_TryPop(t *testing.T) {
	q := NewQueue[int]()
	q.Push(7)
	v, ok := q.TryPop()
	if !ok {
		t.Fatal("TryPop() ok = false, want true")
	}
	if v != 7 {
		t.Errorf("TryPop() = %d, want 7", v)
	}
	if _, ok := q.TryPop(); ok {
		t.Error("TryPop() on drained queue = true, want false")
	}
}

func TestQueue_Close_Push(t *testing.T) {
	q := NewQueue[int]()
	q.Close()
	if ok := q.Push(1); ok {
		t.Error("Push() after Close() = true, want false")
	}
}

func TestQueue_Close_Pop(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	q.Close()
	// Items pushed before Close should still be consumable.
	v, ok := q.Pop()
	if !ok || v != 1 {
		t.Errorf("Pop() = (%d, %v), want (1, true)", v, ok)
	}
	v, ok = q.Pop()
	if !ok || v != 2 {
		t.Errorf("Pop() = (%d, %v), want (2, true)", v, ok)
	}
	if _, ok := q.Pop(); ok {
		t.Error("Pop() on closed empty queue = true, want false")
	}
}

func TestQueue_BlockingPop(t *testing.T) {
	q := NewQueue[int]()
	result := make(chan int, 1)
	go func() {
		v, ok := q.Pop()
		if ok {
			result <- v
		}
	}()
	time.Sleep(10 * time.Millisecond)
	q.Push(42)
	select {
	case v := <-result:
		if v != 42 {
			t.Errorf("blocked Pop() = %d, want 42", v)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked Pop() timed out")
	}
}

func TestQueue_Concurrent(t *testing.T) {
	const (
		producers = 4
		consumers = 4
		itemsEach = 250
	)
	q := NewQueue[int]()
	var wgProd, wgCons sync.WaitGroup
	wgProd.Add(producers)
	for i := 0; i < producers; i++ {
		go func() {
			defer wgProd.Done()
			for j := 0; j < itemsEach; j++ {
				q.Push(1)
			}
		}()
	}
	var (
		mu    sync.Mutex
		total int
	)
	wgCons.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func() {
			defer wgCons.Done()
			for {
				_, ok := q.Pop()
				if !ok {
					return
				}
				mu.Lock()
				total++
				mu.Unlock()
			}
		}()
	}
	wgProd.Wait()
	q.Close()
	wgCons.Wait()
	if total != producers*itemsEach {
		t.Errorf("total consumed = %d, want %d", total, producers*itemsEach)
	}
}

// --- TaskScheduler ---

func TestTaskScheduler_Basic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewTaskScheduler(ctx, 4)
	const n = 20
	results := make([]int, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		if err := s.Go(func() {
			results[i] = i * 2
			wg.Done()
		}); err != nil {
			t.Fatalf("Go() error = %v", err)
		}
	}
	wg.Wait()
	for i := 0; i < n; i++ {
		if results[i] != i*2 {
			t.Errorf("results[%d] = %d, want %d", i, results[i], i*2)
		}
	}
	cancel()
	if err := s.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
}

func TestTaskScheduler_Panic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewTaskScheduler(ctx, 1)
	done := make(chan struct{})
	if err := s.Go(func() {
		defer close(done)
		panic("scheduler panic")
	}); err != nil {
		t.Fatalf("Go() error = %v", err)
	}
	<-done
	cancel()
	err := s.Wait()
	if !isPanicError(err) {
		t.Errorf("Wait() = %v, want ErrPanic", err)
	}
}

func TestTaskScheduler_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewTaskScheduler(ctx, 2)
	cancel()
	if err := s.Wait(); err != nil {
		t.Errorf("Wait() after cancel = %v, want nil", err)
	}
	if err := s.Go(func() {}); !errors.Is(err, ErrClosed) {
		t.Errorf("Go() after cancel = %v, want ErrClosed", err)
	}
}

func TestTaskScheduler_ErrClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewTaskScheduler(ctx, 2)
	s.Close()
	if err := s.Go(func() {}); !errors.Is(err, ErrClosed) {
		t.Errorf("Go() after Close() = %v, want ErrClosed", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
}

func TestTaskScheduler_SingleWorkerOrdering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewTaskScheduler(ctx, 1)
	const n = 100
	results := make([]int, 0, n)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		if err := s.Go(func() {
			mu.Lock()
			results = append(results, i)
			mu.Unlock()
			wg.Done()
		}); err != nil {
			t.Fatalf("Go() error = %v", err)
		}
	}
	wg.Wait()
	for idx, v := range results {
		if v != idx {
			t.Errorf("results[%d] = %d, want %d (FIFO order not maintained)", idx, v, idx)
		}
	}
	cancel()
	_ = s.Wait()
}
