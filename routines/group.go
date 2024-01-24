package routines

import (
	"errors"
	"sync"
)

var (
	ErrClosed           = errors.New("routines group is closed")
	ErrConcurrencyLimit = errors.New("concurrency limit is exceeded")
)

type Group interface {
	Go(func()) error
	Close()
	CloseC() <-chan struct{}
	Wait()
}

type group struct {
	wg      sync.WaitGroup
	closeCh chan struct{}
}

func NewGroup() Group {
	g := &group{
		closeCh: make(chan struct{}),
	}
	return g
}

func (g *group) wrapper(f func()) {
	defer g.wg.Done()
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

func (g *group) Wait() {
	g.wg.Wait()
}

type limitedGroup struct {
	wg        sync.WaitGroup
	routineCh chan struct{}
	closeCh   chan struct{}
}

func NewLimitedGroup(limit uint32) Group {
	g := &limitedGroup{
		routineCh: make(chan struct{}, limit),
		closeCh:   make(chan struct{}),
	}
	return g
}

func (g *limitedGroup) wrapper(f func()) {
	defer func() {
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

func (g *limitedGroup) Wait() {
	g.wg.Wait()
}
