package ratelimit

import (
	"time"
)

type PoolRateLimiter struct {
	ticker *time.Ticker
	pool   chan struct{}
	stop   chan struct{}
}

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

func (p *PoolRateLimiter) Stop() {
	close(p.stop)
}

func (p *PoolRateLimiter) Acquire() {
	<-p.pool
}
