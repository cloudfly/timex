package timex

import (
	"sync/atomic"
	"time"
)

type Limiter struct {
	max  uint64
	c    uint64
	d    time.Duration
	stop chan struct{}
}

func NewLimiter(d time.Duration, c uint64) *Limiter {
	tl := &Limiter{
		max:  c,
		d:    d,
		stop: make(chan struct{}),
	}
	go tl.active()
	return tl
}

func (tl *Limiter) active() {
	ticker := time.NewTicker(tl.d)
	defer ticker.Stop()
	for {
		select {
		case <-tl.stop:
			return
		case <-ticker.C:
			atomic.StoreUint64(&tl.c, 0)
		}
	}
}

func (tl *Limiter) Pass() bool {
	return atomic.AddUint64(&tl.c, 1) <= tl.max
}

func (tl *Limiter) Stop() {
	close(tl.stop)
}
