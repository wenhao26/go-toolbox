package threadpool

// Worker结构体
type Worker struct {
	// 无缓冲的任务队列
	JobQueue JobChan
	// 退出标志
	Quit chan bool
}

// 创建一个新的worker对象
func NewWorker() Worker {
	return Worker{
		make(JobChan),
		make(chan bool),
	}
}

// 启动一个Worker，来监听Job事件
// 执行完成任务，需要将自己重新发送到workerPool
func (w Worker) Start(wp *WorkerPool) {
	// 需要启动一个新的协程，从而不会阻塞
	go func() {
		for {
			// 将Worker注册到线程池
			wp.WorkerQueue <- &w
			select {
			case job := <-w.JobQueue:
				job.RunTask(nil)
			// 终止当前Worker
			case <-w.Quit:
				return
			}
		}
	}()
}
