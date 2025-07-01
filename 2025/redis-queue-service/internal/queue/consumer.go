package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// redisConsumer Consumer接口的Redis具体实现
// 通过Redis的列表数据结构来实现队列功能，使用BRPOP阻塞式地获取消息
type redisConsumer struct {
	client       *redis.Client
	queueName    string
	concurrency  int
	blockTimeout time.Duration
	stopChan     chan struct{} // 用于向所有消费者协程发送停止信号
	wg           sync.WaitGroup
}

// NewRedisConsumer 创建一个新的redisConsumer实例
func NewRedisConsumer(client *redis.Client, queueName string, concurrency int, blockTimeout time.Duration) *redisConsumer {
	return &redisConsumer{
		client:       client,
		queueName:    queueName,
		concurrency:  concurrency,
		blockTimeout: blockTimeout,
		stopChan:     make(chan struct{}),
	}
}

// Consume 启动消费者，根据并发数创建多个协程来并行处理消息
func (rc *redisConsumer) Consume(handler func(message []byte) error) error {
	if rc.concurrency <= 0 {
		return fmt.Errorf("消费者并发性必须大于0 - %d", rc.concurrency)
	}

	rc.wg.Add(rc.concurrency)

	for i := 0; i < rc.concurrency; i++ {
		go func(workerID int) {
			defer rc.wg.Done()

			for {
				select {
				case <-rc.stopChan: // 收到停止信号
					fmt.Printf("消费者[%d]已停止\n", workerID)
					return
				default:
					// BRPOP命令从列表的右侧弹出元素，如果列表为空，则阻塞
					ctx, cancel := context.WithTimeout(context.Background(), rc.blockTimeout+1*time.Second) // BRPOP超时+额外缓冲
					cmd := rc.client.BRPop(ctx, rc.blockTimeout, rc.queueName)
					cancel()

					result, err := cmd.Result()
					if err == redis.Nil { // BRPop超时，没有新消息
						continue
					}
					if err != nil {
						// 对于网络错误等，短暂休眠后重试，避免CPU空转或错误日志刷屏
						time.Sleep(500 * time.Millisecond)
						continue
					}

					// BRPOP 返回一个包含队列名和消息内容的字符串切片
					// result[0] 是队列名，result[1] 是消息内容
					message := []byte(result[1])

					// 调用业务处理函数
					if processErr := handler(message); processErr != nil {
						// 消息处理失败。这里是生产级服务需要重点考虑的地方：
						// 1. **重试机制：** 对于瞬时错误，可以实现指数退避重试
						// 2. **死信队列 (DLQ)：** 将处理失败的消息发送到另一个专门的队列（死信队列），以便后续人工审查或自动处理，避免消息丢失
						// 3. **告警：** 触发告警通知运维人员

						// TODO: 在这里实现死信队列或重试逻辑
						fmt.Println("消息数据:", string(message))
						fmt.Printf("消息处理失败")

					} else {
						fmt.Println("-- OK --")
					}
				}
			}
		}(i)
	}

	rc.wg.Wait()
	return nil
}

// Stop 优雅停止消费者
func (rc *redisConsumer) Stop() {
	close(rc.stopChan) // 关闭通道，所有监听该通道的goroutine将收到信号
}
