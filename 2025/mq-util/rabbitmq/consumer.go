package rabbitmq

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// Consumer 消费者接口
type Consumer interface {
	StartConsuming()
	Close()
}

// consumer 消费者结构体
type consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      string
	handler    func(message []byte) error
	wg         sync.WaitGroup // 等待 goroutine 完成
	shutdownCh chan struct{}  // 用于关闭消费者
}

// NewConsumer 创建RabbitMQ消费者
func NewConsumer(url, queue string, handler func(message []byte) error) *consumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	return &consumer{
		conn:       conn,
		channel:    channel,
		queue:      queue,
		handler:    handler,
		shutdownCh: make(chan struct{}),
	}
}

// StartConsuming 启动消费
func (c *consumer) StartConsuming() {
	msgs, err := c.channel.Consume(c.queue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// 处理消息
	for msg := range msgs {
		select {
		case <-c.shutdownCh:
			// 如果接收到关闭信号，停止消费
			log.Println("Gracefully shutting down consumer.")
			return
		default:
			c.wg.Add(1)
			go func(m amqp.Delivery) {
				defer c.wg.Done()

				// 处理消息
				err := c.handler(m.Body)
				if err != nil {
					log.Printf("Failed to process message: %v", err)
					// 如果处理失败，可以根据需求实现重试机制
				} else {
					log.Printf("Processed message: %s", string(m.Body))
				}
			}(msg)
		}
	}
}

// Close 关闭消费者
func (c *consumer) Close() {
	// 关闭消费者时，发送关闭信号并等待处理完所有消息
	close(c.shutdownCh)
	c.wg.Wait() // 等待所有消息处理完

	_ = c.channel.Close()
	_ = c.conn.Close()
}
