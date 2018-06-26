package concurrent

import "sync/atomic"

func (worker *Worker) run() {
	go func() {
		for run := range worker.task {
			if run == nil {
				atomic.AddInt32(&worker.pool.going, -1)
				return
			}
			run()
			worker.pool.addWorker(worker)
		}
	}()
}

func (worker *Worker) interrupt() {
	worker.addTask(nil)
}

func (worker *Worker) addTask(task Run) {
	worker.task <- task
}
