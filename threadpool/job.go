package threadpool

// 需要执行的job
type Job interface {
	RunTask(request interface{})
}

// job channel
type JobChan chan Job
