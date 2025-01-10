package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// 消费者结构体
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *Config
}

// 创建消费者实例
func NewConsumer(config *Config) (*Consumer, error) {
	conn, err := amqp.Dial(config.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ连接失败: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("打开channel失败: %w", err)
	}

	// 设置Qos
	err = ch.Qos(config.RabbitMQ.PrefetchCount, 0, false)
	if err != nil {
		_ = conn.Close()
		_ = ch.Close()
		return nil, fmt.Errorf("设置QoS失败: %w", err)
	}

	return &Consumer{conn: conn, channel: ch, config: config}, nil
}

// 启动消费者
func (c *Consumer) Start() {
	var wg sync.WaitGroup

	queueName := c.config.RabbitMQ.QueueName

	// 为每个消费者启动协程
	for i := 0; i < c.config.RabbitMQ.ConsumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()

			msgs, err := c.channel.Consume(
				queueName,                                 // 队列名称
				fmt.Sprintf("go_consumer_%d", consumerID), // 消费者标签
				false,                                     // 是否自动确认消息
				false,
				false,
				false, // 是否阻塞
				nil,
			)
			if err != nil {
				log.Fatalf("Consumer %d，未能注册消费者: %s", consumerID, err)
			}

			workerPool := make(chan bool, c.config.RabbitMQ.WorkerPoolSize)

			for msg := range msgs {
				workerPool <- true

				// 使用协程处理消息
				go func(message amqp.Delivery) {
					defer func() { <-workerPool }()
					c.handleMessage(consumerID, message)
				}(msg)
			}

		}(i + 1)
	}

	wg.Wait()
}

// 处理消息的业务逻辑
func (c *Consumer) handleMessage(consumerID int, msg amqp.Delivery) {
	log.Printf("Consumer %d，收到的消息: %s", consumerID, msg.Body)

	// 模拟业务逻辑
	if err := processMessage(string(msg.Body)); err != nil {
		log.Printf("Consumer %d: Failed to process message: %s", consumerID, err)
		// 消息处理失败，重新投递到队列
		_ = msg.Nack(false, true)
		return
	}
	// 消息处理成功，确认
	_ = msg.Ack(false)
}

// 模拟处理消息的业务逻辑
func processMessage(body string) error {
	log.Printf("Processing message: %s", body)
	// 在这里写实际的业务逻辑
	return nil
}

// 关闭RabbitMQ相关连接
func (c *Consumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
}
