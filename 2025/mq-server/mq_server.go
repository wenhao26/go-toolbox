// 实现MQ服务器，处理生产者和消费者的连接和消息传递
package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"toolbox/2025/mq-server/manager"
)

// MQServer 消息队列服务器
type MQServer struct {
	addr        string
	queueMgr    *manager.QueueManager
	connections map[*net.Conn]struct{}
	lock        sync.Mutex
}

// NewMQServer 创建一个新的消息队列服务器
func NewMQServer(addr string) *MQServer {
	return &MQServer{
		addr:        addr,
		queueMgr:    manager.NewQueueManager(),
		connections: make(map[*net.Conn]struct{}),
	}
}

// Start 启动服务器
func (s *MQServer) Start() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
	log.Printf("MQ服务器启动，监听地址: %s", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}
		s.lock.Lock()
		s.connections[&conn] = struct{}{}
		s.lock.Unlock()
		go s.handleConnection(conn)
	}
}

// handleConnection 处理客户端连接
func (s *MQServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.lock.Lock()
		delete(s.connections, &conn)
		s.lock.Unlock()
	}()

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	for {
		var request struct {
			Type      string `gob:"type"`
			QueueName string `gob:"queueName"`
			Message   string `gob:"message"`
		}

		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("解码请求失败: %v", err)
			continue
		}

		switch request.Type {
		case "produce":
			queue := s.queueMgr.GetQueue(request.QueueName, 100)
			_, _ = queue.Publish(request.Message)
			_ = encoder.Encode("消息已发送")
		case "consume":
			queue := s.queueMgr.GetQueue(request.QueueName, 100)
			message := queue.Receive()
			if message != "" {
				_ = encoder.Encode(message)
			} else {
				_ = encoder.Encode("")
			}
		default:
			log.Printf("未知请求类型: %s", request.Type)
			_ = encoder.Encode("未知请求类型")
		}
	}
}

// Close 关闭服务器
func (s *MQServer) Close() {
	s.lock.Lock()
	for conn := range s.connections {
		conn.Close()
	}
	s.queueMgr.CloseAllQueues()
	s.lock.Unlock()
}

func main() {
	server := NewMQServer(":8080")
	defer server.Close()
	server.Start()
}
