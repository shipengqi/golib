package lock

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

/*
Test Result:
		goos: windows
		goarch: amd64
		pkg: github.com/shipengqi/golib/safe/lock
		cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
		BenchmarkMutex-12               16091377                69.01 ns/op
		BenchmarkSpinLock-12            50514362                23.04 ns/op
		BenchmarkBackOffSpinLock-12     63326876                19.27 ns/op

*/

type spinLockWithoutEb uint32

func (l *spinLockWithoutEb) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(l), 0, 1) {
		runtime.Gosched()
	}
}

func (l *spinLockWithoutEb) Unlock() {
	atomic.StoreUint32((*uint32)(l), 0)
}

func newSpinLockWithoutEb() sync.Locker {
	return new(spinLockWithoutEb)
}

func BenchmarkMutex(b *testing.B) {
	m := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			m.Unlock()
		}
	})
}

func BenchmarkSpinLock(b *testing.B) {
	spin := newSpinLockWithoutEb()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			spin.Unlock()
		}
	})
}

func BenchmarkBackOffSpinLock(b *testing.B) {
	spin := NewSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			spin.Unlock()
		}
	})
}
