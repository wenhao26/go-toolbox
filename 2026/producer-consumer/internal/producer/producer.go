// Package producer 实现了模拟“实时数据推送”的生产者
package producer

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"toolbox/2026/producer-consumer/internal/model"
	"toolbox/2026/producer-consumer/internal/stats"
)

// Producer 模拟一个独立的实时数据源
// 系统中可以并行运行多个 Producer 实例，共同向同一个 channel 写入数据
// 从而模拟真实业务中“多路数据汇聚”的高并发写入场景
type Producer struct {
	id     int                   // 生产者编号，用于日志追踪与数据溯源
	qps    int                   // 目标生产速率（每秒消息数）
	out    chan<- *model.Message // 只写 channel 生产者只能写入，不能读取，编译期保证方向安全
	seq    *atomic.Uint64        // 全局唯一自增序列号生成器，由多个 Producer 共享以保证 ID 全局唯一
	stats  *stats.Stats
	logger *slog.Logger
}

// New 创建一个 Producer 实例
// seq 由调用方（App）统一持有并在多个 Producer 之间共享
// 因为 atomic.Uint64 本身就是并发安全的计数器，无需额外加锁
func New(id, qps int, out chan<- *model.Message, seq *atomic.Uint64, stats *stats.Stats, logger *slog.Logger) *Producer {
	return &Producer{
		id:     id,
		qps:    qps,
		out:    out,
		seq:    seq,
		stats:  stats,
		logger: logger.With("component", "producer", "producer_id", id),
	}
}

// Run 持续生产消息知道 ctx 被取消，这是个阻塞方法，通常配合 go p.Run(ctx, wg) 使用
//
// 限速实现：使用 time.Ticker 按固定间隔触发生产，是 Go 中实现简单限流最常见、开销最低的方式
// 相比 sleep 循环，Ticker 有 runtime 统一调度，误差更小
func (p *Producer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if p.qps <= 0 {
		p.qps = 1
	}

	interval := time.Second / time.Duration(p.qps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	p.logger.Info("生产者已启动", "qps", p.qps)
	defer p.logger.Info("生产者已停止")

	for {
		select {
		case <-ctx.Done():
			// 收到停止信号：不再生产新消息，直接退出
			return
		case <-ticker.C:
			msg := p.buildMessage()
			// 关键点：写入 channel 时，同样监听 ctx.Done()
			// 避免在队列已满（消费能力跟不上）且系统正在关闭时被永久阻塞
			// 这是生产者优雅退出不可或缺的一环
			select {
			case p.out <- msg:
				p.stats.IncProduced()
			case <-ctx.Done():
				return
			}
		}
	}
}

// buildMessage 构造一条模拟消息，随机分配消息类型，模拟真实业务中的负载混合特征
func (p *Producer) buildMessage() *model.Message {
	id := p.seq.Add(1)
	return &model.Message{
		ID:         id,
		ProducerID: p.id,
		Type:       randomMessageType(),
		Payload:    fmt.Sprintf("payload-from-producer-%d-seq-%d", p.id, id),
		CreatedAt:  time.Now(),
	}
}

// randomMessageType 按权重随机生成消息类型：60% 轻量 | 30% 中等 | 10% 重量级
// 贴近真实系统“大部分请求轻量，少部分请求昂贵”的长尾分布特征
func randomMessageType() model.MessageType {
	r := rand.Float64()
	switch {
	case r < 0.6:
		return model.TypeLight
	case r < 0.9:
		return model.TypeMedium
	default:
		return model.TypeHeavy
	}
}
