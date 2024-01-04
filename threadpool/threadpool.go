package threadpool

import (
	"sync"
)

var workerPool *WorkerPool
var once sync.Once

// 线程池
type WorkerPool struct {
	Size        int
	JobQueue    JobChan
	WorkerQueue chan *Worker
}

// 单例模式来获取workerPool
func GetWorkerPool(poolSize, jobQueueLen int) *WorkerPool {
	once.Do(func() {
		workerPool = NewWorkerPool(poolSize, jobQueueLen)
	})
	return workerPool
}

func NewWorkerPool(poolSize, jobQueueLen int) *WorkerPool {
	return &WorkerPool{
		poolSize,
		make(JobChan, jobQueueLen),
		make(chan *Worker, poolSize),
	}
}

func (wp *WorkerPool) Start() {

	// 将所有worker启动
	for i := 0; i < wp.Size; i++ {
		worker := NewWorker()
		worker.Start(wp)
	}

	// 监听JobQueue，如果接收到请求，随机取一个Worker，然后将Job发送给该Worker的JobQueue
	// 需要启动一个新的协程，来保证不阻塞
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker.JobQueue <- job
			}
		}
	}()
}
