package safe

import (
	"context"
	"log"
	"sync"
)

type RecoverCallback func(err error)

// DefaultRecoverCallback .
func DefaultRecoverCallback(err error) {
	buf := stack(2)
	log.Printf("safe go panic\t%v\n%s", err, buf)
}

func Go(handler func(), options ...Options) {
	go func() {
		defer recovery(options...)
		handler()
	}()
}

func GoAndWait(handlers ...func() error) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)

	for _, fn := range handlers {
		wg.Add(1)
		go func(handler func() error) {
			defer recovery()
			defer wgdone(&wg)

			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(fn)
	}
	wg.Wait()
	return err
}

func GoAndWaitWithCtx(ctx context.Context, handlers ...func(ctx context.Context) error) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)

	for _, fn := range handlers {
		wg.Add(1)
		go func(handler func(ctx context.Context) error) {
			defer recovery()
			defer wgdone(&wg)

			if e := handler(ctx); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(fn)
	}
	wg.Wait()
	return err
}

func recovery(options ...Options) {
	if e := recover(); e != nil {
		buf := stack(2)
		log.Printf("safe go panic\t%v\n%s", e, buf)
	}
}

func wgdone(wg *sync.WaitGroup) {
	wg.Done()
}
