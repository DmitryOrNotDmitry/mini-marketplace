package ratelimit

import (
	"time"
)

// PoolRateLimiter реализует ограничение RPS с помощью пула с периодическим пополнением.
type PoolRateLimiter struct {
	ticker *time.Ticker
	pool   chan struct{}
	stop   chan struct{}
}

// NewPoolRateLimiter создает новый PoolRateLimiter.
func NewPoolRateLimiter(poolSize int, tickerDuration time.Duration) *PoolRateLimiter {
	p := &PoolRateLimiter{
		ticker: time.NewTicker(tickerDuration),
		pool:   make(chan struct{}, poolSize),
		stop:   make(chan struct{}),
	}

	go func() {
		defer p.ticker.Stop()

		for {
			select {
			case <-p.ticker.C:
				select {
				case p.pool <- struct{}{}:
					continue
				case <-p.stop:
					return
				}
			case <-p.stop:
				return
			}
		}
	}()

	return p
}

// Stop останавливает работу PoolRateLimiter.
func (p *PoolRateLimiter) Stop() {
	close(p.stop)
}

// Acquire блокирует до получения разрешения из пула.
func (p *PoolRateLimiter) Acquire() {
	<-p.pool
}
