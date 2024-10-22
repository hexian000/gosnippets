// gosnippets (c) 2023-2024 He Xian <hexian000@outlook.com>
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

type ErrPanic struct {
	v any
}

func (p ErrPanic) Panic() any {
	return p.v
}

func (p ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", p.v)
}

var _ = error(ErrPanic{})

type Group interface {
	Go(func()) error
	Close()
	CloseC() <-chan struct{}
	Wait() error
}

type group struct {
	wg      sync.WaitGroup
	closeCh chan struct{}
	errorCh chan error
}

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
