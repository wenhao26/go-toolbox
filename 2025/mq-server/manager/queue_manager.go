// 定义消息队列管理器，用于管理多个队列
package manager

import (
	"errors"
	"sync"
)

// MessageQueue 消息队列结构体
type MessageQueue struct {
	name    string      // 名称
	channel chan string // 存放消息管道
	mu      sync.Mutex  // 并发锁
}

// NewMessageQueue 创建一个消息队列的实例
func NewMessageQueue(name string, bufferSize int) *MessageQueue {
	return &MessageQueue{
		name:    name,
		channel: make(chan string, bufferSize),
	}
}

// Publish 向队列推送消息
func (mq *MessageQueue) Publish(message string) (bool, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// 非阻塞写入
	select {
	case mq.channel <- message:
		return true, nil
	default:
		return false, errors.New("消息写入失败，管道缓冲区已满")
	}
}

// Receive 从队列取出消息
func (mq *MessageQueue) Receive() string {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	select {
	case message := <-mq.channel:
		return message
	default:
		return ""
	}
}

// Close 关闭队列
func (mq *MessageQueue) Close() {
	close(mq.channel)
}

// QueueManager 消息队列管理器结构体
type QueueManager struct {
	queues map[string]*MessageQueue
	mu     sync.Mutex
}

// NewQueueManager 创建一个新的消息队列管理器
func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[string]*MessageQueue),
	}
}

// GetQueue 获取或创建消息队列
func (qm *QueueManager) GetQueue(name string, bufferSize int) *MessageQueue {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if qm.queues[name] == nil {
		qm.queues[name] = NewMessageQueue(name, bufferSize)
	}
	return qm.queues[name]
}

// CloseQueue 关闭队列
func (qm *QueueManager) CloseQueue(name string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if q, ok := qm.queues[name]; ok {
		q.Close()
		delete(qm.queues, name)
	}
}

// CloseAllQueues 关闭所有队列
func (qm *QueueManager) CloseAllQueues() {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	for _, q := range qm.queues {
		q.Close()
	}
	qm.queues = make(map[string]*MessageQueue)
}
