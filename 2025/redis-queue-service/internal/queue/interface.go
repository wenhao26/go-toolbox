package queue

// Producer 定义向队列发送消息的通用接口
type Producer interface {
	// Enqueue 将消息字节数组发送到队列
	// 成功返回nil，否则返回相关错误
	Enqueue(message []byte) error
}

// Consumer 定义从队列接收和处理消息的通用接口
type Consumer interface {
	// Consume 启动消费者，从队列中获取消息并执行处理函数
	// - handler:处理函数，用于处理从队列取出的消息
	// 如果消费过程中发生不可恢复的错误，将返回相关错误
	Consume(handler func(message []byte) error) error

	// Stop 优雅停止消费者
	// 它会等待所有正在处理的消费完成，并关闭相关资源
	Stop()
}
