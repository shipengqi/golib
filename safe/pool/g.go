package pool

import (
	"fmt"
	"runtime"
)

type G struct {
	// pool who owns this G.
	pool *Pool
}

func (g *G) run() {
	g.pool.addRunning(1)
	go func() {
		defer func() {
			g.pool.addRunning(-1)
			g.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := g.pool.panicHandler; ph != nil {
					ph(p)
				} else {
					// Todo custom Logger
					fmt.Printf("worker exits from a panic: %v\n", p)
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					fmt.Printf("worker exits from panic: %s\n", string(buf[:n]))
				}
			}
			// Call Signal() here in case there are goroutines waiting for available workers.
			g.pool.cond.Signal()
		}()

		for f := range g.task {
			if f == nil {
				return
			}
			f()
			if ok := g.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
