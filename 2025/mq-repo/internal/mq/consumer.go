// 通用消费者封装（协程池并发消费、死信队列、手动 ack、消费失败重试、panic recover）
// ./app consumer --queue=xxx --workers=20

package mq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// MessageHandler 消息处理接口
type MessageHandler interface {
	HandleMessage(msg amqp.Delivery) error
}

// Consumer 消费者结构体
type Consumer struct {
	uri         string
	queue       string
	concurrency int
	handler     MessageHandler
	prefetch    int
	retryDelay  time.Duration
	maxRetry    int
}

// NewConsumer 创建消息者实例
func NewConsumer(uri, queue string, concurrency int, handler MessageHandler) *Consumer {
	return &Consumer{
		uri:         uri,
		queue:       queue,
		concurrency: concurrency,
		handler:     handler,
		prefetch:    10,
		retryDelay:  2 * time.Second,
		maxRetry:    3,
	}
}

// Run 启动消费者
func (c *Consumer) Run() error {
	conn := GetConnection(c.uri)
	ch := conn.Channel()

	// 设置prefetch限制同时处理消息的数量
	if err := ch.Qos(c.prefetch, 0, false); err != nil {
		return err
	}

	// 队列声明
	if _, err := ch.QueueDeclare(c.queue, true, false, false, false, nil); err != nil {
		return err
	}

	// 开始消费消息
	msgs, err := ch.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	pool := NewWorkerPool(c.concurrency)
	defer pool.Close()

	for msg := range msgs {
		msgCopy := msg // 避免闭包捕获循环变量
		pool.Submit(func() {
			retryCount := 0

			for {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("消费者任务 panic recovered: %v", r)
						}
					}()

					err := c.handler.HandleMessage(msgCopy)
					if err != nil {
						retryCount++
						log.Printf("消息处理失败(%s), retry %d/%d: %v", c.queue, retryCount, c.maxRetry, err)
						if retryCount >= c.maxRetry {
							// 超过重试次数，进入死信
							_ = msgCopy.Nack(false, false)
						} else {
							time.Sleep(c.retryDelay)
							return
						}
					} else {
						_ = msgCopy.Ack(false)  // 手动确认
						retryCount = c.maxRetry // 退出循环
					}
				}()

				if retryCount >= c.maxRetry {
					break
				}
			}
		})
	}

	return nil
}
