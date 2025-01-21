// 实现消费者客户端，用于从MQ服务器获取消息

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

// Consumer 消费者
type Consumer struct {
	conn net.Conn
}

// NewConsumer 创建一个新的消费者
func NewConsumer(addr string) (*Consumer, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Consumer{conn: conn}, nil
}

// Consume 获取消息
func (c *Consumer) Consume(queueName string) (string, error) {
	encoder := gob.NewEncoder(c.conn)
	request := struct {
		Type      string `gob:"type"`
		QueueName string `gob:"queueName"`
	}{
		Type:      "consume",
		QueueName: queueName,
	}
	if err := encoder.Encode(request); err != nil {
		return "", err
	}
	var message string
	decoder := gob.NewDecoder(c.conn)
	if err := decoder.Decode(&message); err != nil {
		return "", err
	}
	return message, nil
}

// Close 关闭连接
func (c *Consumer) Close() error {
	return c.conn.Close()
}

func main() {
	// 创建消费者
	consumer, err := NewConsumer("localhost:8080")
	if err != nil {
		log.Fatalf("创建消费者失败: %v", err)
	}
	defer consumer.Close()

	// 消费者获取消息
	for {
		message, err := consumer.Consume("testQueue")
		if err != nil {
			log.Fatalf("获取消息失败: %v", err)
		}
		if message == "" {
			fmt.Println("没有消息可消费")
		} else {
			fmt.Printf("消费消息: %s\n", message)
		}
		time.Sleep(1 * time.Second)
	}
}
