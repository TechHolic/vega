package vega

import (
	"sync"
	"runtime"
	"testing"
	"time"
	"github.com/techholic/vega/concurrent"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (超过了int32的范围)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424    (超过了int64的范围)
	YiB // 1208925819614629174706176
)

var n = 1000000

func demoFunc() error {
	n := 100
	time.Sleep(time.Duration(n) * time.Millisecond)
	return nil
}

func TestCachedPool(t *testing.T) {
	var wg sync.WaitGroup
	cachedRoutinePool := concurrent.CachedRoutinePool()
	for i := 0; i < n; i++ {
		wg.Add(1)
		cachedRoutinePool.Submit(func() error {
			demoFunc()
			wg.Done()
			return nil
		})
	}
	wg.Wait()
	t.Logf("going workers number:%d", cachedRoutinePool.GetGoing())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	t.Logf("memory usage:%d MB", mem.TotalAlloc/MiB)
}

func TestFixPool(t *testing.T) {
	var wg sync.WaitGroup
	fixRoutinePool := concurrent.FixRoutinePool(int32(200000))
	for i := 0; i < n; i++ {
		wg.Add(1)
		fixRoutinePool.Submit(func() error {
			demoFunc()
			wg.Done()
			return nil
		})
	}
	wg.Wait()
	t.Logf("going workers number:%d", fixRoutinePool.GetGoing())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	t.Logf("memory usage:%d MB", mem.TotalAlloc/MiB)
}

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc()
			wg.Done()
		}()
	}
	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	t.Logf("memory usage:%d MB", mem.TotalAlloc/MiB)
}
