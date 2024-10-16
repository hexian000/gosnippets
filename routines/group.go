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
	Go(func() error) error
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

func (g *group) wrapper(f func() error) {
	defer g.wg.Done()
	if err := f(); err != nil {
		select {
		case g.errorCh <- err:
		default:
		}
	}
}

func (g *group) Go(f func() error) error {
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

func (g *limitedGroup) wrapper(f func() error) {
	defer func() {
		<-g.routineCh
		g.wg.Done()
	}()
	if err := f(); err != nil {
		select {
		case g.errorCh <- err:
		default:
		}
	}
}

func (g *limitedGroup) Go(f func() error) error {
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
