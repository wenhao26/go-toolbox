// 实现生产者客户端，用于发送消息到MQ服务器

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

// Producer 生产者
type Producer struct {
	conn net.Conn
}

// NewProducer 创建一个新的生产者
func NewProducer(addr string) (*Producer, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Producer{conn: conn}, nil
}

// Produce 发送消息
func (p *Producer) Produce(queueName, message string) error {
	encoder := gob.NewEncoder(p.conn)
	request := struct {
		Type      string `gob:"type"`
		QueueName string `gob:"queueName"`
		Message   string `gob:"message"`
	}{
		Type:      "produce",
		QueueName: queueName,
		Message:   message,
	}
	if err := encoder.Encode(request); err != nil {
		return err
	}
	var response string
	decoder := gob.NewDecoder(p.conn)
	if err := decoder.Decode(&response); err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

// Close 关闭连接
func (p *Producer) Close() error {
	return p.conn.Close()
}

func main() {
	// 创建生产者
	producer, err := NewProducer("localhost:8080")
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}
	defer producer.Close()

	// 生产者发送消息
	if err := producer.Produce("testQueue", "Hello, World!"); err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}
}
