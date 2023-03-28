package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shipengqi/golib/safe/lock"
)

var (
	ErrInvalidPoolCapacity = errors.New("")
	ErrInvalidPoolValidity = errors.New("")
)

// Pool accept the tasks from client,it limits the total
// of goroutines to a given number by recycling goroutines.
type Pool struct {
	// capacity of the pool.
	capacity int32

	// validity sets up the period of validity (second) for the G.
	validity time.Duration

	// panicHandler sets up the panic handler.
	panicHandler func(interface{})

	// waiting is the number of goroutines already been blocked on pool.Go(), protected by pool.lock
	waiting int32

	// idle is a slice that store the available G.
	idle []*G

	// running is the number of the running goroutines.
	running int32

	// cond for waiting to get an idle worker.
	cond *sync.Cond

	// release is used to notice all goroutines to close.
	release chan struct{}

	// lock for protecting the G queue.
	lock sync.Locker
}

// NewPool create a new goroutine Pool.
func NewPool(capacity int, options ...Option) (*Pool, error) {
	p := &Pool{
		capacity: int32(capacity),
		lock:     lock.NewSpinLock(),
	}
	p.withOptions(options...)

	if p.capacity <= 0 {
		return nil, ErrInvalidPoolCapacity
	}
	if p.validity <= 0 {
		return nil, ErrInvalidPoolValidity
	}

	p.cond = sync.NewCond(p.lock)

	return p, nil
}

// Running returns the number of G currently running.
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Waiting returns the number of tasks which are waiting be executed.
func (p *Pool) Waiting() int {
	return int(atomic.LoadInt32(&p.waiting))
}

// Cap returns the capacity of this pool.
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Idle returns the number of idle G, -1 indicates this pool is unlimited.
func (p *Pool) Idle() int {
	c := p.Cap()
	if c < 0 {
		return -1
	}
	return c - p.Running()
}

// Tune changes the capacity of this pool, note that it is noneffective to the infinite or pre-allocation pool.
func (p *Pool) Tune(capacity int) {
	current := p.Cap()
	if current == -1 || capacity <= 0 || capacity == current {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(capacity))
	if capacity > current {
		if capacity-current == 1 {
			p.cond.Signal()
			return
		}
		p.cond.Broadcast()
	}
}

// withOptions set options for the application
func (p *Pool) withOptions(opts ...Option) *Pool {
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

func (p *Pool) addRunning(num int) {
	atomic.AddInt32(&p.running, int32(num))
}

func (p *Pool) addWaiting(num int) {
	atomic.AddInt32(&p.waiting, int32(num))
}
