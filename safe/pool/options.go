package pool

import (
	"time"
)

type Option interface {
	apply(p *Pool)
}

// optionFunc wraps a func, so it satisfies the Option interface.
type optionFunc func(*Pool)

func (f optionFunc) apply(p *Pool) {
	f(p)
}

// WithPanicHandler sets up the panic handler.
func WithPanicHandler(panicHandler func(interface{})) Option {
	return optionFunc(func(p *Pool) {
		p.panicHandler = panicHandler
	})
}

// WithValidity sets up the period of validity (second) for the G.
func WithValidity(validity time.Duration) Option {
	return optionFunc(func(p *Pool) {
		p.validity = validity
	})
}
