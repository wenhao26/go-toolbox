package main

import (
	"fmt"
	"log"

	"toolbox/2025/mq-util/rabbitmq2"
)

func main() {
	// 初始化连接池，连接到RabbitMQ服务，创建5个连接池
	rabbitURL := "amqp://admin:admin@localhost:5672/"
	pool, err := rabbitmq2.NewConnectionPool(rabbitURL, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close() // 程序结束时关闭所有连接

	// 创建生产者
	producerConn := pool.GetConnection() // 从连接池中获取一个连接
	producer, err := rabbitmq2.NewProducer(producerConn, "testExchange", "testQueue")
	if err != nil {
		log.Fatal(err)
	}

	// 发送消息
	err = producer.SendMessage("Hello RabbitMQ with Exchange!")
	if err != nil {
		log.Fatal(err)
	}

	// 创建消费者
	consumerConn := pool.GetConnection() // 从连接池中获取另一个连接
	consumer, err := rabbitmq2.NewConsumer(consumerConn, "testExchange", "testQueue")
	if err != nil {
		log.Fatal(err)
	}

	// 消费消息并打印
	err = consumer.ConsumeMessages(func(msg string) {
		fmt.Println("Received message:", msg)
	})
	if err != nil {
		log.Fatal(err)
	}
}
