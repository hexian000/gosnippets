// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package routines

import (
	"context"
	"sync"
)

// TaskScheduler represents a scheduler of tasks with limited parallelism.
type TaskScheduler struct {
	wg      sync.WaitGroup
	queue   *Queue[func()]
	errorCh chan error
}

// NewTaskScheduler creates a task scheduler that runs at most numWorkers
// tasks concurrently. When ctx is cancelled, the scheduler is closed and
// remaining queued tasks are discarded.
func NewTaskScheduler(ctx context.Context, numWorkers int) *TaskScheduler {
	s := &TaskScheduler{
		queue:   NewQueue[func()](),
		errorCh: make(chan error, 1),
	}
	s.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go s.worker()
	}
	go s.watchCtx(ctx)
	return s
}

func (s *TaskScheduler) worker() {
	defer s.wg.Done()
	for {
		f, ok := s.queue.Pop()
		if !ok {
			return
		}
		s.run(f)
	}
}

func (s *TaskScheduler) run(f func()) {
	defer func() {
		if v := recover(); v != nil {
			select {
			case s.errorCh <- &ErrPanic{v}:
			default:
			}
		}
	}()
	f()
}

func (s *TaskScheduler) watchCtx(ctx context.Context) {
	<-ctx.Done()
	s.Close()
}

// Go enqueues a task. Returns ErrClosed if the scheduler has been closed.
func (s *TaskScheduler) Go(f func()) error {
	if !s.queue.Push(f) {
		return ErrClosed
	}
	return nil
}

// Close closes the task scheduler. Queued but unstarted tasks are discarded.
func (s *TaskScheduler) Close() {
	s.queue.Close()
}

// Wait waits for all running tasks to finish and returns the first
// panic error, if any.
func (s *TaskScheduler) Wait() error {
	s.wg.Wait()
	select {
	case err := <-s.errorCh:
		return err
	default:
	}
	return nil
}
