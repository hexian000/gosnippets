// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package routines

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrClosed           = errors.New("routines group is closed")
	ErrConcurrencyLimit = errors.New("concurrency limit is exceeded")
)

// ErrPanic represents a panic error from a goroutine.
type ErrPanic struct {
	v any
}

// Panic returns the value passed to panic().
func (p ErrPanic) Panic() any {
	return p.v
}

// Error implements the error interface.
func (p ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", p.v)
}

var _ = error(ErrPanic{})

// Group represents a group of goroutines.
type Group interface {
	// Go starts a new goroutine in the group.
	Go(func()) error
	// Close signals that no more goroutines will be started.
	Close()
	// CloseC returns a channel that is closed when the group is closed.
	// It can be used to cancel long-running goroutines.
	CloseC() <-chan struct{}
	// Wait waits for all goroutines in the group to finish.
	Wait() error
}

type group struct {
	wg      sync.WaitGroup
	closeCh chan struct{}
	errorCh chan error
}

// NewGroup creates and returns a new Group.
func NewGroup() Group {
	g := &group{
		closeCh: make(chan struct{}),
		errorCh: make(chan error, 1),
	}
	return g
}

func (g *group) wrapper(f func()) {
	defer func() {
		if v := recover(); v != nil {
			select {
			case g.errorCh <- ErrPanic{v}:
			default:
			}
		}
		g.wg.Done()
	}()
	f()
}

func (g *group) Go(f func()) error {
	select {
	case <-g.closeCh:
		return ErrClosed
	default:
	}
	g.wg.Add(1)
	go g.wrapper(f)
	return nil
}

func (g *group) Close() {
	close(g.closeCh)
}

func (g *group) CloseC() <-chan struct{} {
	return g.closeCh
}

func (g *group) Wait() error {
	g.wg.Wait()
	select {
	case err := <-g.errorCh:
		return err
	default:
	}
	return nil
}

type limitedGroup struct {
	wg        sync.WaitGroup
	routineCh chan struct{}
	closeCh   chan struct{}
	errorCh   chan error
}

// NewLimitedGroup creates and returns a new Group with a concurrency limit.
func NewLimitedGroup(limit int) Group {
	g := &limitedGroup{
		routineCh: make(chan struct{}, limit),
		closeCh:   make(chan struct{}),
		errorCh:   make(chan error, 1),
	}
	return g
}

func (g *limitedGroup) wrapper(f func()) {
	defer func() {
		if v := recover(); v != nil {
			select {
			case g.errorCh <- &ErrPanic{v}:
			default:
			}
		}
		<-g.routineCh
		g.wg.Done()
	}()
	f()
}

func (g *limitedGroup) Go(f func()) error {
	select {
	case <-g.closeCh:
		return ErrClosed
	case g.routineCh <- struct{}{}:
	default:
		return ErrConcurrencyLimit
	}
	g.wg.Add(1)
	go g.wrapper(f)
	return nil
}

func (g *limitedGroup) Close() {
	close(g.closeCh)
}

func (g *limitedGroup) CloseC() <-chan struct{} {
	return g.closeCh
}

func (g *limitedGroup) Wait() error {
	g.wg.Wait()
	select {
	case err := <-g.errorCh:
		return err
	default:
	}
	return nil
}
