// Package consumer 实现了固定大小的 worker pool，用于并发消费消息队列
package consumer

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"
	"toolbox/2026/producer-consumer/internal/model"
	"toolbox/2026/producer-consumer/internal/stats"
)

// costRange 定义了不同消息类型对应的模拟处理耗时区间 [min, max]
var costRange = map[model.MessageType][2]time.Duration{
	model.TypeLight:  {5 * time.Millisecond, 20 * time.Millisecond},
	model.TypeMedium: {50 * time.Millisecond, 120 * time.Millisecond},
	model.TypeHeavy:  {200 * time.Millisecond, 500 * time.Millisecond},
}

// Pool 是消费端的核心：一组共享同一输入 channel 的 worker goroutine
// Go 的 channel 是多读多写并发安全的，因此多个 worker 同时对同一个 channel 执行 range/receive，
// 天然实现了“竞争消费、多劳多得”的负载均衡，无需额外的分发器或锁
type Pool struct {
	workerNum int
	in        <-chan *model.Message // 只读 channel：worker 只能读取，不能写入
	stats     *stats.Stats
	logger    *slog.Logger
}

// NewPool 创建一个拥有 workerNum 个并发 worker 的消费池
func NewPool(workerNum int, in <-chan *model.Message, stat *stats.Stats, logger *slog.Logger) *Pool {
	return &Pool{
		workerNum: workerNum,
		in:        in,
		stats:     stat,
		logger:    logger.With("component", "consumer"),
	}
}

// Run 启动 workerNum 个 goroutine 开始消费，非阻塞方法，立即返回
// 每个 worker 通过 wg 向上层汇报退出，供 App 编排优雅停机流程
func (p *Pool) Run(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < p.workerNum; i++ {
		wg.Add(1)
		go p.worker(ctx, i, wg)
	}
	p.logger.Info("工作者池已启动", "worker_num", p.workerNum)
}

// worker 是单个消费协程的主循环
//
// 设计要点（优雅退出核心）：
// 1、主循环用 “for msg := range p.in” 而不是 select + ctx.Done
// 这样即便收到了停机信号，只要 channel 尚未关闭且里面还有消息，
// worker 仍会继续把队列中已经生产好的消息处理完，不会“腰斩”业务数据
//
// 2、channel 何时关闭由上层 App 统一控制：先停止所有生产者，
// 再 close(channel)，worker 的 range 循环便会在消费完剩余元素后自然退出
//
// 3、ctx 只用于“强制超时兜底”：如果优雅退出耗时超时预设阈值，上层会 cancel 这个 ctx，
// 正在 process() 中等待的 worker 会被立刻打断
func (p *Pool) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// panic 兜底：单条消息处理逻辑出现 panic 不应导致整个 worker 崩溃
	// 这是生产级消费者必须具备的健壮性设计
	defer p.recoverWorkerPanic(id)

	p.logger.Debug("工作者已启动", "worker_id", id)

	for msg := range p.in {
		p.process(ctx, id, msg)
	}

	p.logger.Debug("工作者已退出（输入通道已关闭）", "worker_id", id)
}

// recoverWorkerPanic 恢复工作进程
func (p *Pool) recoverWorkerPanic(id int) {
	if r := recover(); r != nil {
		p.logger.Error("工作者协程死机已恢复", "worker_id", id, "panic", r)
	}
}

// process 处理单条消息，并模拟与消息类型相关的处理耗时
func (p *Pool) process(ctx context.Context, workerID int, msg *model.Message) {
	// 单条消息级别的 panic 兜底：即使某条消息的业务逻辑异常，
	// 也只影响这一条消息，worker 主循环继续处理下一条
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error("消息处理死机已恢复",
				"worker_id", workerID,
				"msg_id", msg.ID,
				"panic", r,
			)
			p.stats.IncFailed()
		}
	}()

	p.stats.IncConsumed()

	cost := randomCost(msg.Type)
	timer := time.NewTimer(cost)
	defer timer.Stop()

	select {
	case <-timer.C:
		// 正常处理完成
		p.stats.IncSucceeded()
		p.logger.Info("消息已处理",
			"worker_id", workerID,
			"msg_id", msg.ID,
			"type", msg.Type.String(),
			"process_cost", cost,
			"e2e_latency", time.Since(msg.CreatedAt),
		)
	case <-ctx.Done():
		// 只有在触发了“强制超时兜底”时才会走到这个分支：
		// 当前消息处理被中断，计入失败，供运维观察优雅退出是否触发了强制截断
		p.stats.IncFailed()
		p.logger.Warn("消息处理因强制关机中断",
			"worker_id", workerID,
			"msg_id", msg.ID,
		)
	}

}

// randomCost 在指定消息类型的耗时区间内返回一个随机时长，模拟真实处理耗时的抖动
func randomCost(t model.MessageType) time.Duration {
	rng, ok := costRange[t]
	if !ok {
		rng = costRange[model.TypeLight]
	}
	lo, hi := rng[0], rng[1]
	span := int64(hi - lo)
	if span <= 0 {
		return lo
	}
	return lo + time.Duration(rand.Int64N(span+1))
}
