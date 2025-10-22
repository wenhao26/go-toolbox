// RabbitMQ连接封装（生产级单例连接池，自动重连，心跳检测）

package mq

import (
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// Connection 封装了RabbitMQ连接和channel
type Connection struct {
	uri      string // RabbitMQ连接地址
	conn     *amqp.Connection
	channel  *amqp.Channel
	mu       sync.Mutex // 保护conn和channel并发安全
	isClosed bool
}

var (
	globalConn *Connection
	once       sync.Once
)

// GetConnection 获取全局单例RabbitMQ连接
func GetConnection(uri string) *Connection {
	once.Do(func() {
		globalConn = &Connection{uri: uri}
		if err := globalConn.connect(); err != nil {
			log.Fatalf("RabbitMQ初始化失败：%v", err)
		}
	})
	return globalConn
}

// connect 建立连接并创建channel，支持自动重连
func (c *Connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 支持心跳和超时
	conn, err := amqp.DialConfig(c.uri, amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	})
	if err != nil {
		return err
	}

	// 创建channel
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	c.conn = conn
	c.channel = ch
	c.isClosed = false

	// 异步监控连接断开并重连
	go c.handleReconnect()

	return nil
}

// handleReconnect 监听连接断开事件，自动重连
func (c *Connection) handleReconnect() {
	errChan := make(chan *amqp.Error)
	c.conn.NotifyClose(errChan)

	for {
		err := <-errChan
		if err != nil {
			log.Printf("RabbitMQ连接断开，尝试重连：%v\n", err)

			for {
				time.Sleep(5 * time.Second) // 重连间隔
				if err := c.connect(); err == nil {
					log.Println("RabbitMQ重连成功")
					return
				}
				log.Println("RabbitMQ重连失败，继续重试...")
			}
		}
	}
}

// Channel 获取当前可用channel
func (c *Connection) Channel() *amqp.Channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.channel
}

// Close 优雅关闭RabbitMQ连接和channel
func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isClosed {
		_ = c.channel.Close()
		_ = c.conn.Close()
		c.isClosed = true
	}
}
