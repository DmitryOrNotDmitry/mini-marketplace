package myerrgroup

import (
	"context"
	"sync"
	"sync/atomic"
)

type ErrorGroup struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	err    error
	wasErr atomic.Bool
}

func New() *ErrorGroup {
	return &ErrorGroup{}
}

func WithContext(ctx context.Context) (*ErrorGroup, context.Context) {
	derivedCtx, cancel := context.WithCancel(ctx)

	return &ErrorGroup{
		cancel: cancel,
	}, derivedCtx
}

func (g *ErrorGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if g.wasErr.Load() {
			return
		}

		err := f()
		if err != nil {
			if g.wasErr.CompareAndSwap(false, true) {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			}
		}
	}()
}

func (g *ErrorGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
