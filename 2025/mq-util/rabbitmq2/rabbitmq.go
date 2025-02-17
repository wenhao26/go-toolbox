// 管理RabbitMQ连接和交换机/队列的绑定
package rabbitmq2

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// ConnectionPool 管理连接池
// 连接池用于生产和消费过程中复用RabbitMQ连接，避免每次都重新创建连接，提高性能
type ConnectionPool struct {
	Connections []*amqp.Connection
}

// NewConnectionPool 初始化连接池
// 根据给定的URL和池的大小，创建多个连接，并将其存储在连接池中
func NewConnectionPool(rabbitURL string, poolSize int) (*ConnectionPool, error) {
	var pool []*amqp.Connection

	// 创建多个连接
	for i := 0; i < poolSize; i++ {
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
		}
		pool = append(pool, conn)
	}

	return &ConnectionPool{Connections: pool}, nil
}

// GetConnection 从连接池中获取一个连接
// 每次从连接池中返回一个连接，池中的连接会在使用后被取出
// 如果连接池为空，这里返回nil（实际场景中可以考虑阻塞或重试逻辑）
func (cp *ConnectionPool) GetConnection() *amqp.Connection {
	if len(cp.Connections) == 0 {
		log.Fatal("No available connections in the pool.")
		return nil
	}

	// 取出并返回第一个连接
	conn := cp.Connections[0]
	cp.Connections = cp.Connections[1:]

	return conn
}

// Close 关闭所有连接
// 清理连接池中的所有连接
func (cp *ConnectionPool) Close() {
	for _, conn := range cp.Connections {
		_ = conn.Close() // 关闭连接
	}
}
