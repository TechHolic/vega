package concurrent

import (
	"math"
	"runtime"
)

func CachedRoutinePool() (*RoutinePool) {
	cachedRoutinePool, _ := New(math.MaxInt32)
	return cachedRoutinePool
}

func FixRoutinePool(size int32) (*RoutinePool) {
	if size <= 0 {
		size = int32(runtime.NumCPU())
	}
	fixRoutinePool, _ := New(size)
	return fixRoutinePool
}
