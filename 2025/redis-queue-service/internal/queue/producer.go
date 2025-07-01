package queue

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// redisProducer Producer接口的Redis具体实现
// 通过Redis的列表数据结构来实现队列功能，使用LPUSH将消息推入
type redisProducer struct {
	client    *redis.Client
	queueName string
}

// NewRedisProducer 创建一个新的redisProducer实例
func NewRedisProducer(client *redis.Client, queueName string) *redisProducer {
	return &redisProducer{
		client:    client,
		queueName: queueName,
	}
}

// Enqueue 实现Producer接口的Enqueue方法
// 将消息推送到redis列表的头部
func (rp *redisProducer) Enqueue(message []byte) error {
	ctx := context.Background()

	// LPUSH命名将一个或多个值插入到列表头部
	if err := rp.client.LPush(ctx, rp.queueName, message).Err(); err != nil {
		return fmt.Errorf("消息推送到队列失败 %s: %w", rp.queueName, err)
	}

	return nil
}
