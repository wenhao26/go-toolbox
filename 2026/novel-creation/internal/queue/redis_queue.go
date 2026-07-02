package queue

import (
	"context"
	"errors"
	"fmt"
	"time"
	"toolbox/2026/novel-creation/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisQueue struct {
	client          *redis.Client
	sourceQueue     string // 源队列名，如：creation:task
	processingQueue string // 处理中队列名，如：creation:task:processing
	blockTimeout    time.Duration
}

func NewRedisQueue(client *redis.Client, sourceQueue, processingSuffix string, blockTimeout int) *RedisQueue {
	return &RedisQueue{
		client:          client,
		sourceQueue:     sourceQueue,
		processingQueue: processingSuffix,
		blockTimeout:    time.Duration(blockTimeout) * time.Second,
	}
}

// PopTask 原子地将任务从源队列移动到处理中队列，返回任务内容
func (q *RedisQueue) PopTask(ctx context.Context) (string, error) {
	val, err := q.client.BRPopLPush(ctx, q.sourceQueue, q.processingQueue, q.blockTimeout).Result()
	if errors.Is(err, redis.Nil) {
		// 超时无任务，正常情况
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("BRPopLPush failed: %w", err)
	}
	return val, nil
}

// AckTask 任务处理成功后，从处理中队列删除该任务
func (q *RedisQueue) AckTask(ctx context.Context, task string) error {
	cnt, err := q.client.LRem(ctx, q.processingQueue, 1, task).Result()
	if err != nil {
		return fmt.Errorf("LRem failed: %w", err)
	}
	if cnt == 0 {
		logger.Log.Warn("task not found in processing queue", zap.String("task", task))
	}
	return nil
}

// RequeueTask 任务处理失败时，将其放回源队列（可选，也可记录到死信队列）
func (q *RedisQueue) RequeueTask(ctx context.Context, task string) error {
	// 先从处理中队列删除，再推入源队列
	if err := q.AckTask(ctx, task); err != nil {
		return err
	}
	return q.client.LPush(ctx, q.sourceQueue, task).Err()
}

// RecoverOrphanTasks 程序启动时，将处理中队列所有的任务，重新返回源队列，避免任务丢失
func (q *RedisQueue) RecoverOrphanTasks(ctx context.Context) error {
	for {
		val, err := q.client.RPopLPush(ctx, q.processingQueue, q.sourceQueue).Result()
		if errors.Is(err, redis.Nil) {
			break
		}
		if err != nil {
			return err
		}
		logger.Log.Info("recovered orphan task", zap.String("task", val))
	}
	return nil
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}
