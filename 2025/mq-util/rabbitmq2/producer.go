// 生产者封装，支持交换机绑定和消息发布
package rabbitmq2

import (
	"github.com/streadway/amqp"
)

// Producer 生产者相关操作
// 生产者负责连接到RabbitMQ，声明交换机与队列，发送消息
type Producer struct {
	Channel   *amqp.Channel // AMQP通道，用于发送和接收消息
	Queue     amqp.Queue    // 队列对象，存储消息
	Exchange  string        // 交换机名称，用于消息路由
	QueueName string        // 队列名称
}

// NewProducer 创建一个新的生产者
// 通过连接获取AMQP通道，声明交换机和队列，并绑定交换机到队列
func NewProducer(conn *amqp.Connection, exchange, queueName string) (*Producer, error) {
	// 获取通道
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 声明交换机（如果交换机不存在则创建）
	err = channel.ExchangeDeclare(
		exchange, // 交换机名称
		"direct", // 交换机类型：direct, topic, fanout等
		true,     // 是否持久化交换机
		false,    // 是否自动删除交换机
		false,    // 是否独占
		false,    // 是否设置其他属性
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 声明队列（如果队列不存在则创建）
	queue, err := channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化队列
		false,     // 是否自动删除队列
		false,     // 是否独占
		false,     // 是否阻塞
		nil,       // 其他参数
	)
	if err != nil {
		return nil, err
	}

	// 绑定交换机到队列
	err = channel.QueueBind(
		queue.Name, // 队列名
		"",         // 路由键，空字符串表示所有消息
		exchange,   // 交换机名称
		false,      // 是否阻塞
		nil,        // 其他参数
	)
	if err != nil {
		return nil, err
	}

	// 返回生产者对象
	return &Producer{
		Channel:   channel,
		Queue:     queue,
		Exchange:  exchange,
		QueueName: queueName,
	}, nil
}

// SendMessage 发送消息到队列
// 将消息发布到指定的交换机，经过路由键被投递到队列中
func (p *Producer) SendMessage(message string) error {
	// 发布消息到交换机
	err := p.Channel.Publish(
		p.Exchange, // 交换机名称
		"",         // 路由键（空字符串表示所有消息）
		false,      // 是否持久化
		false,      // 是否立即发送
		amqp.Publishing{
			ContentType: "text/plain",    // 消息类型
			Body:        []byte(message), // 消息内容
		},
	)

	return err
}
