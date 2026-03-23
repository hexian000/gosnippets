// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package routines

import "sync"

// Queue is a MPMC (Multiple Producer Multiple Consumer) unbounded queue
// that uses a swap-buffer technique for high throughput. Two slices are
// maintained: producers append to the write slice while consumers read
// from the read slice. When the read slice is drained, the two slices
// are swapped so the old read buffer is reused for future writes,
// minimizing allocations and lock contention per item.
type Queue[T any] struct {
	mu      sync.Mutex
	cond    *sync.Cond
	write   []T
	read    []T
	readOff int
	closed  bool
}

// NewQueue creates a new unbounded MPMC queue.
func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push enqueues a value. Returns false if the queue is closed.
func (q *Queue[T]) Push(v T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.write = append(q.write, v)
	q.cond.Signal()
	return true
}

// Pop dequeues a value, blocking until one is available.
// Returns false if the queue is closed and empty.
func (q *Queue[T]) Pop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for {
		if q.readOff < len(q.read) {
			v := q.read[q.readOff]
			var zero T
			q.read[q.readOff] = zero
			q.readOff++
			return v, true
		}
		// read buffer is drained, try to swap
		if len(q.write) > 0 {
			q.read, q.write = q.write, q.read[:0]
			q.readOff = 0
			continue
		}
		if q.closed {
			var zero T
			return zero, false
		}
		q.cond.Wait()
	}
}

// TryPop tries to dequeue a value without blocking.
// Returns false if the queue is empty or closed.
func (q *Queue[T]) TryPop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.readOff < len(q.read) {
		v := q.read[q.readOff]
		var zero T
		q.read[q.readOff] = zero
		q.readOff++
		return v, true
	}
	if len(q.write) > 0 {
		q.read, q.write = q.write, q.read[:0]
		q.readOff = 0
		v := q.read[0]
		var zero T
		q.read[0] = zero
		q.readOff = 1
		return v, true
	}
	var zero T
	return zero, false
}

// Len returns the number of items currently in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.write) + len(q.read) - q.readOff
}

// Close closes the queue. After closing, Push returns false and
// Pop returns false once all remaining items have been consumed.
func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
