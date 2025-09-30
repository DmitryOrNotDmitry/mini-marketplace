package myerrgroup

import (
	"context"
	"sync"
	"sync/atomic"
)

// ErrorGroup реализует группу горутин с обработкой ошибок и отменой контекста.
type ErrorGroup struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	err    error
	wasErr atomic.Bool
}

// New создает новый ErrorGroup.
func New() *ErrorGroup {
	return &ErrorGroup{}
}

// WithContext создает ErrorGroup с поддержкой отмены контекста.
func WithContext(ctx context.Context) (*ErrorGroup, context.Context) {
	derivedCtx, cancel := context.WithCancel(ctx)

	return &ErrorGroup{
		cancel: cancel,
	}, derivedCtx
}

// Go запускает функцию в отдельной горутине, отслеживая ошибку.
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

// Wait ожидает завершения всех горутин и возвращает первую ошибку.
func (g *ErrorGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
