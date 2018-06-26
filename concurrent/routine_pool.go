package concurrent

import (
	"errors"
	"math"
	"sync/atomic"
)

var (
	PoolSizeInvalidError = errors.New("invalid size for pool")
	PoolClosedError      = errors.New("this pool has been closed")
)

func New(size int32) (*RoutinePool, error) {
	if size <= 0 {
		return nil, PoolSizeInvalidError
	}
	p := &RoutinePool{
		capacity:      int32(size),
		freeSignal:    make(chan Signal, math.MaxInt32),
		releaseSignal: make(chan Signal, 1),
	}
	return p, nil
}

func (pool *RoutinePool) Submit(task Run) error {
	if len(pool.releaseSignal) > 0 {
		return PoolClosedError
	}
	worker := pool.getWorker()
	worker.addTask(task)
	return nil
}

func (pool *RoutinePool) GetGoing() int32 {
	return atomic.LoadInt32(&pool.going)
}

func (pool *RoutinePool) GetFree() int32 {
	return atomic.LoadInt32(&pool.capacity) - atomic.LoadInt32(&pool.going)
}

func (pool *RoutinePool) GetCapacity() int32 {
	return atomic.LoadInt32(&pool.capacity)
}

func (pool *RoutinePool) Resize(size int32) {
	if size < pool.capacity {
		diff := size - pool.capacity
		for i := int32(0); i < diff; i++ {
			pool.getWorker().interrupt()
		}
	} else if size == pool.capacity {
		return
	}
	atomic.StoreInt32(&pool.capacity, size)
}

func (pool *RoutinePool) Close() error {
	pool.once.Do(func() {
		pool.releaseSignal <- Signal{}
		going := pool.going
		for i := int32(0); i < going; i++ {
			pool.getWorker().interrupt()
		}
		for i := range pool.workers {
			pool.workers[i] = nil
		}
	})
	return nil
}

func (pool *RoutinePool) getWorker() *Worker {
	var worker *Worker
	waiting := false
	pool.lock.Lock()
	workers := pool.workers
	n := len(workers) - 1
	if n < 0 {
		if pool.going >= pool.capacity {
			waiting = true
		} else {
			pool.going++
		}
	} else {
		<-pool.freeSignal
		worker = workers[n]
		workers[n] = nil
		pool.workers = workers[:n]
	}
	pool.lock.Unlock()
	if waiting {
		<-pool.freeSignal
		pool.lock.Lock()
		workers = pool.workers
		l := len(workers) - 1
		worker = workers[l]
		workers[l] = nil
		pool.workers = workers[:l]
		pool.lock.Unlock()
	} else if worker == nil {
		worker = &Worker{
			pool: pool,
			task: make(chan Run),
		}
		worker.run()
	}
	return worker
}

func (pool *RoutinePool) addWorker(worker *Worker) {
	pool.lock.Lock()
	pool.workers = append(pool.workers, worker)
	pool.lock.Unlock()
	pool.freeSignal <- Signal{}
}
