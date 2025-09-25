package myerrgroup

import (
	"context"
	"sync"
)

type ErrorGroup struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	err    error
	errMx  sync.Mutex
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

		err := f()
		if err != nil && g.err == nil {
			g.errMx.Lock()
			if g.err == nil {
				g.err = err

				if g.cancel != nil {
					g.cancel()
				}
			}
			g.errMx.Unlock()
		}
	}()
}

func (g *ErrorGroup) Wait() error {
	g.wg.Wait()
	return g.err
}
