// 消费者封装，支持交换机绑定和消费消息
package rabbitmq2

import (
	"github.com/streadway/amqp"
)

// Consumer 消费者相关操作
// 消费者负责从队列中接收消息，并根据业务逻辑处理
type Consumer struct {
	Channel   *amqp.Channel // AMQP通道，用于发送和接收消息
	Queue     amqp.Queue    // 队列对象，存储消息
	Exchange  string        // 交换机名称，用于消息路由
	QueueName string        // 队列名称
}

// NewConsumer 创建一个新的消费者
// 获取通道，声明交换机和队列，并绑定交换机到队列
func NewConsumer(conn *amqp.Connection, exchange, queueName string) (*Consumer, error) {
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

	// 返回消费者对象
	return &Consumer{
		Channel:   channel,
		Queue:     queue,
		Exchange:  exchange,
		QueueName: queueName,
	}, nil
}

// ConsumeMessages 消费消息并处理
func (c *Consumer) ConsumeMessages(handler func(string)) error {
	// 从队列中消费消息
	msgs, err := c.Channel.Consume(
		c.Queue.Name, // 队列名称
		"",           // 消费者标签
		true,         // 是否自动确认
		false,        // 是否独占
		false,        // 是否阻塞
		false,        // 是否取消
		nil,          // 其他参数
	)
	if err != nil {
		return err
	}

	// 处理消息
	for msg := range msgs {
		handler(string(msg.Body)) // 调用业务处理函数
	}

	return nil
}
