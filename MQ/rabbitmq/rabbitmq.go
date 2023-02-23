package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const Url = "amqp://test:test@127.0.0.1:5672/"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	QueueName string
	Exchange  string
	Key       string

	MqUrl string
}

var err error

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     Url,
	}
}

func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) FailErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("err=%s", err))
	}
}

// 简单模式生产者
func (r *RabbitMQ) PublishSimple(message string) {
	// 声明队列
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,
	)
	r.FailErr(err)

	// 调用channel发送消息到队列
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	log.Printf("Sent=%s OK!", message)
}

// 简单模式消费者
func (r *RabbitMQ) ConsumerSimple() {
	// 声明队列
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 接收消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			// todo 消息逻辑处理
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 简单模式
func NewSimple(queueName string) *RabbitMQ {
	rabbitMq := NewRabbitMQ(queueName, "", "")
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.FailErr(err)

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.FailErr(err)

	return rabbitMq
}

// 订阅模式生产者
func (r *RabbitMQ) PublishPub(message string) {
	// 声明交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 发送消息
	r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	log.Printf("Sent=%s OK!", message)
}

// 订阅模式消费者
func (r *RabbitMQ) ReceiveSub() {
	// 声明队列
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 声明队列，注意这里队列名称不要填写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.FailErr(err)

	// 绑定队列到exchange中
	r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	// 接收消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	forever := make(chan bool)
	go func() {
		for message := range messages {
			// todo 消息逻辑处理
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 订阅模式
func NewSub(exchange string) *RabbitMQ {
	rabbitMq := NewRabbitMQ("", exchange, "")
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.FailErr(err)

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.FailErr(err)

	return rabbitMq
}

// 路由模式生产者
func (r *RabbitMQ) RoutingPub(message string) {
	// 声明交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 发送消息
	r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	log.Printf("Sent=%s OK!", message)
}

func (r *RabbitMQ) ReceiveRouting() {
	// 声明队列
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 声明队列，注意这里队列名称不要填写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.FailErr(err)

	// 绑定队列到exchange中
	r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	// 接收消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	forever := make(chan bool)
	go func() {
		for message := range messages {
			// todo 消息逻辑处理
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 路由模式
func NewRouting(exchange, routing string) *RabbitMQ {
	rabbitMq := NewRabbitMQ("", exchange, routing)
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.FailErr(err)

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.FailErr(err)

	return rabbitMq
}

// 话题模式生产者
func (r *RabbitMQ) TopicPub(message string) {
	// 声明交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 发送消息
	r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	log.Printf("Sent=%s OK!", message)
}

// 话题模式消费者
// 要注意key,规则
// 其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
// 匹配 test.* 表示匹配 test.hello, test.hello.one需要用test.#才能匹配到
func (r *RabbitMQ) TopicConsumer() {
	// 声明队列
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	// 声明队列，注意这里队列名称不要填写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.FailErr(err)

	// 绑定队列到exchange中
	r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	// 接收消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.FailErr(err)

	forever := make(chan bool)
	go func() {
		for message := range messages {
			// todo 消息逻辑处理
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 话题模式
func NewTopic(exchange, routing string) *RabbitMQ {
	rabbitMq := NewRabbitMQ("", exchange, routing)
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.FailErr(err)

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.FailErr(err)

	return rabbitMq
}
