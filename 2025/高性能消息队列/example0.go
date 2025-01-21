package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// 需求：基于 Go Channel 的高性能消息队列。
//功能点：
// - 使用 Channel 实现消息队列。
// - 支持多生产者和多消费者。
// - 实现消息的持久化（可选）。

// MessageQueue 消息队列结构体
type MessageQueue struct {
	queue chan string
	wg    sync.WaitGroup
}

// NewMessageQueue 创建消息队列实例
func NewMessageQueue(bufferSize int) *MessageQueue {
	return &MessageQueue{
		queue: make(chan string, bufferSize),
	}
}

// Produce 生产
func (m *MessageQueue) Produce(message string) {
	m.queue <- message
}

// Consume 消费
func (m *MessageQueue) Consume() {
	defer m.wg.Done()
	for message := range m.queue {
		log.Println(message)
		time.Sleep(100 * time.Millisecond)
	}
}

// Close 关闭
func (m *MessageQueue) Close() {
	close(m.queue)
	m.wg.Wait()
}

// SystemInfo 运行系统信息
func (m *MessageQueue) SystemInfo() string {
	// CPU 使用率
	cpus, _ := cpu.Percent(1*time.Second, false)

	// 内存使用情况
	memInfo, _ := mem.VirtualMemory()

	return fmt.Sprintf(
		"CPU使用率: %.2f%%\n"+
			"总内存: %.2f GB\n"+
			"已用内存: %.2f GB\n"+
			"内存使用率: %.2f%%\n",
		cpus[0],
		float64(memInfo.Total)/1024/1024/1024,
		float64(memInfo.Used)/1024/1024/1024,
		memInfo.UsedPercent,
	)
}

func main() {
	var count int
	mq := NewMessageQueue(10)

	// 启动消费者
	mq.wg.Add(1)
	go mq.Consume()

	// 生产者模拟推送数据
	for {
		if count > 100000 {
			mq.Close()
			return
		}
		message := mq.SystemInfo()
		count++
		mq.Produce(message)
	}
}
