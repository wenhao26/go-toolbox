package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Producer 生产者接口
type Producer interface {
	Publish(message string) error
	Close()
}

// producer 生产者结构体
type producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// NewProducer 创建RabbitMQ生产者
func NewProducer(url, queue string) *producer {
	conn, _ := amqp.Dial(url)
	channel, _ := conn.Channel()
	_, _ = channel.QueueDeclare(queue, true, false, false, false, nil)

	return &producer{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}
}

// Publish 发送消息
func (p *producer) Publish(message string) error {
	return p.channel.Publish("", p.queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

// Close 关闭生产者
func (p *producer) Close() {
	_ = p.channel.Close()
	_ = p.conn.Close()
}
