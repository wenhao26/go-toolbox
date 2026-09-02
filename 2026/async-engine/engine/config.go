package engine

// Config 定义异步引擎的运行参数
type Config struct {
	WorkerNum int // 并发消费协程数量
	QueueSize int // 任务队列缓冲长度
}
