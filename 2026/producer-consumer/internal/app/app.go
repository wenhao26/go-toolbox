// Package app 负责将 producer、consumer、stats 等组件
// 编排为一个完整的、具备优雅停止能力的应用生命周期
package app

import (
	"context"
	"log/slog"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"toolbox/2026/producer-consumer/internal/config"
	"toolbox/2026/producer-consumer/internal/consumer"
	"toolbox/2026/producer-consumer/internal/model"
	"toolbox/2026/producer-consumer/internal/producer"
	"toolbox/2026/producer-consumer/internal/stats"
)

// App 整个生产者-消费者系统的顶层编排者，持有所有共享资源的引用
type App struct {
	cfg    *config.Config
	logger *slog.Logger
	stats  *stats.Stats
	queue  chan *model.Message // 生产者与消费者之间的有界缓冲队列（背压核心）
}

// New 构造 App
// 此时尚未创建 channel、尚未启动任何协程
func New(cfg *config.Config, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		stats:  stats.New(),
	}
}

// Run 是应用主入口，按四个阶段驱动整个生命周期：
//
// 【启动】：创建有界 channel 作为队列，启动 N 个生产者与 M 个消费者 worker
//
// 【运行】：阻塞等待，直到外部传入的 ctx 被取消（收到 SIGINT/Ctrl+C）或达到配置的运行时长
//
// 【优雅停机】：先停止生产者生产新数据 -> 等到生产者协程全部退出 ->
// 关闭 channel（通知消费者“不会再有新数据了”） -> 等待所有 worker 把 channel 中剩余的消息消费完
//
// 【强制兜底】：如果【优雅停止】超过 ShutdownTimeout 仍未完成（例如消息处理逻辑卡死），
// 则强制取消 worker 的处理 ctx，避免进程无法退出
//
// 模型标准范式：先停止生产，再消费完队列，最后设置一个安全阈（超时强制退出）
func (a *App) Run(ctx context.Context) error {
	a.queue = make(chan *model.Message, a.cfg.QueueSize)
	var seq atomic.Uint64

	// producerCtx 用于【优雅停机】主动停止生产，与外部传入的 ctx 解耦
	// 这样即使外部 ctx 已经 Done，我们依然能精确控制“先停止生产、后停止消费”的顺序
	producerCtx, stopProducing := context.WithCancel(context.Background())
	defer stopProducing()

	// workerCtx，仅用于【强制兜底】，正常流程不会被触发
	workerCtx, forceStopWorkers := context.WithCancel(context.Background())
	defer forceStopWorkers()

	// 阶段一：启动
	var producerWG sync.WaitGroup
	for i := 0; i < a.cfg.ProducerNum; i++ {
		producerWG.Add(1)
		p := producer.New(i, a.cfg.ProduceQPS, a.queue, &seq, a.stats, a.logger)
		go p.Run(producerCtx, &producerWG)
	}

	pool := consumer.NewPool(a.cfg.WorkerNum, a.queue, a.stats, a.logger)
	var workerWG sync.WaitGroup
	pool.Run(workerCtx, &workerWG)

	statsCtx, stopStats := context.WithCancel(context.Background())
	defer stopStats()
	go a.reportStats(statsCtx)

	a.logger.Info("系统正在运行",
		"producer_num", a.cfg.ProducerNum,
		"worker_num", a.cfg.WorkerNum,
		"queue_size", a.cfg.QueueSize,
		"qps_per_producer", a.cfg.ProduceQPS,
	)

	// 阶段二：运行
	if a.cfg.RunDuration > 0 {
		select {
		case <-ctx.Done():
			a.logger.Info("收到中断信号")
		case <-time.After(a.cfg.RunDuration):
			a.logger.Info("已达到配置运行时长")
		}
	} else {
		<-ctx.Done()
		a.logger.Info("收到中断信号")
	}

	// 阶段三：优雅停机
	a.logger.Info("优雅关闭：步骤1/3 - 停止生成新消息")
	stopProducing()
	producerWG.Wait()

	a.logger.Info("优雅关闭：步骤2/3 - 关闭队列，清空剩余消息",
		"remaining_in_queue", len(a.queue))
	close(a.queue)

	a.logger.Info("优雅关机：第3步/共3步 - 等待工作进程完成",
		"timeout", a.cfg.ShutdownTimeout)
	workersDone := make(chan struct{})
	go func() {
		workerWG.Wait()
		close(workersDone)
	}()

	// 阶段四：强制兜底
	select {
	case <-workersDone:
		a.logger.Info("所有工作进程均正常退出，未触发超时")
	case <-time.After(a.cfg.ShutdownTimeout):
		a.logger.Warn("关闭超时，强制工作进程终止正在处理的消息")
		forceStopWorkers()
		<-workersDone // 强制取消后，process() 内部的 select 会立即返回，这里很快就能等到
	}

	stopStats()
	a.printFinalStats()
	return nil
}

// reportStats 周期性打印吞吐、队列深度、goroutine 数量等运行时指标，
// 便于在压测过程中直观观察系统的处理能力与积压情况。
func (a *App) reportStats(ctx context.Context) {
	ticker := time.NewTicker(a.cfg.StatsInterval)
	defer ticker.Stop()

	// lastProduced/lastConsumed/lastTime 用于计算相邻两次采样之间的增量速率，
	// 只在 reportStats 这一个 goroutine 里读写，不涉及并发，无需加锁。
	var lastProduced, lastConsumed uint64
	lastTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			produced, consumed, succeeded, failed := a.stats.Snapshot()
			elapsed := now.Sub(lastTime).Seconds()

			var produceRate, consumeRate float64
			if elapsed > 0 {
				produceRate = float64(produced-lastProduced) / elapsed
				consumeRate = float64(consumed-lastConsumed) / elapsed
			}

			a.logger.Info("统计快照",
				"produced", produced,
				"consumed", consumed,
				"succeeded", succeeded,
				"failed", failed,
				"queue_len", len(a.queue),
				"queue_cap", cap(a.queue),
				"produce_rate_per_sec", math.Round(produceRate*10)/10,
				"consume_rate_per_sec", math.Round(consumeRate*10)/10,
				"goroutines", runtime.NumGoroutine(),
			)
		}
	}
}

func (a *App) printFinalStats() {
	produced, consumed, succeeded, failed := a.stats.Snapshot()
	a.logger.Info("最终统计",
		"produced", produced,
		"consumed", consumed,
		"succeeded", succeeded,
		"failed", failed,
	)
}
