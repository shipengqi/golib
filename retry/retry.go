// Package retry provide the retry function.
package retry

import (
	"context"
	"time"
)

type Retry struct {
	ctx      context.Context
	interval time.Duration
	times    int
}

// Times creates a Retry with times.
func Times(times int) *Retry {
	r := &Retry{}
	return r.Times(times)
}

// Times sets the retry times
func (r *Retry) Times(times int) *Retry {
	r.times = times
	return r
}

// WithInterval sets the retry interval.
func (r *Retry) WithInterval(interval time.Duration) *Retry {
	r.interval = interval
	return r
}

// WithContext sets the retry context.
func (r *Retry) WithContext(ctx context.Context) *Retry {
	r.ctx = ctx
	return r
}

// Do execute the given func with retry.
func (r *Retry) Do(f func() error) error {
	var err error
	for ; r.times > 0; r.times-- {
		if r.isDone() {
			break
		}
		if err = f(); err != nil {
			r.sleep()
		} else {
			break
		}
	}
	return err
}

func (r *Retry) isDone() bool {
	if r.ctx == nil {
		return false
	}
	select {
	case <-r.ctx.Done():
		return true
	default:
		return false
	}
}

func (r *Retry) sleep() {
	if r.interval == 0 {
		return
	}
	if r.ctx == nil {
		<-time.After(r.interval)
		return
	}
	select {
	case <-time.After(r.interval):
	case <-r.ctx.Done():
	}
}
