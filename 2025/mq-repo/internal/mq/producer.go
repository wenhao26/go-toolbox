// RabbitMQ生产者封装（异步/同步发布、Confirm模式、失败重试、队列自动声明）
// ./app producer --queue=xxx

package mq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Producer 封装生产者逻辑
type Producer struct {
	ch       *amqp.Channel          // 通道
	confirm  chan amqp.Confirmation // Confirm模式通知
	asyncBuf chan publishTask       // 异步发送缓冲
}

// publishTask 异步发布任务结构体
type publishTask struct {
	queue string
	body  []byte
}

// NewProducer 创建生产者实例
// async: 是否使用异步发布
// bufSize: 异步缓冲队列大小
func NewProducer(uri string, async bool, bufSize int) *Producer {
	conn := GetConnection(uri)
	ch := conn.Channel()

	// 开启confirm模式，保证消息投递可靠
	if err := ch.Confirm(false); err != nil {
		log.Fatalf("启用confirm模式失败：%v", err)
	}

	p := &Producer{
		ch:       ch,
		confirm:  make(chan amqp.Confirmation, 100),
		asyncBuf: make(chan publishTask, bufSize),
	}

	// 异步worker
	if async {
		p.RunAsyncWorker()
	}

	return p
}

// RunAsyncWorker 启动异步worker
func (p *Producer) RunAsyncWorker() {
	go func() {
		for task := range p.asyncBuf {
			if err := p.Publish(task.queue, task.body); err != nil {
				log.Printf("异步发送消息失败：%v", err)
			}
		}
	}()
}

// Publish 同步发布消息，支持confirm模式和失败重试
func (p *Producer) Publish(queue string, body []byte) error {
	// 队列声明，保证队列存在
	if _, err := p.ch.QueueDeclare(
		queue, // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 无等待
		nil,   // 参数
	); err != nil {
		return err
	}

	for i := 0; i < 3; i++ { // 重试3次
		err := p.ch.Publish(
			"",    // 默认交换机
			queue, // routing key
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				Body:         body,
			},
		)
		if err == nil {
			return nil
		}
		log.Printf("Publish失败，重试 %d/3: %v", i+1, err)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// PublishAsync 异步发送消息
func (p *Producer) PublishAsync(queue string, body []byte) {
	p.asyncBuf <- publishTask{queue: queue, body: body}
}

// Close 关闭生产者
func (p *Producer) Close() {
	close(p.asyncBuf)
	_ = p.ch.Close()
}
