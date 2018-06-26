package concurrent

import "sync"

type Signal struct {
}

type Run func() error

type RoutinePool struct {
	capacity int32

	going int32

	freeSignal chan Signal

	workers []*Worker

	releaseSignal chan Signal

	lock sync.Mutex

	once sync.Once
}

type Worker struct {
	pool *RoutinePool

	task chan Run
}
