package stats

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"toolbox/2026/async-engine/engine"
	"toolbox/2026/async-engine/producer"
)

// Reporter 负责定时采集引擎与生产者指标，输出实时看板与最终汇总报表。
type Reporter struct {
	engine   *engine.Engine
	producer *producer.Producer
	interval time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu      sync.Mutex
	history []Record
}

// Record 保存一次采样的完整数据。
type Record struct {
	Time           time.Time
	Engine         engine.Snapshot
	Producer       producer.Snapshot
	InstantQPS     float64
	InstantDrop    int64
	InstantConsume int64
}

// NewReporter 创建报表实例，interval 控制采样频率。
func NewReporter(eng *engine.Engine, prod *producer.Producer, interval time.Duration) *Reporter {
	ctx, cancel := context.WithCancel(context.Background())
	return &Reporter{
		engine:   eng,
		producer: prod,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动后台定时采样协程。
func (r *Reporter) Start() {
	r.wg.Add(1)
	go r.loop()
}

// Stop 停止采样，追加最终快照并输出汇总报表。
func (r *Reporter) Stop() {
	r.cancel()
	r.wg.Wait()

	// 追加最终快照，确保报表包含最新状态
	r.mu.Lock()
	r.history = append(r.history, Record{
		Time:     time.Now(),
		Engine:   r.engine.Snapshot(),
		Producer: r.producer.Snapshot(),
	})
	r.mu.Unlock()

	r.renderFinalReport()
}

func (r *Reporter) loop() {
	defer r.wg.Done()

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	var lastEng engine.Snapshot
	var lastProd producer.Snapshot

	for {
		select {
		case <-r.ctx.Done():
			return
		case now := <-ticker.C:
			engSnap := r.engine.Snapshot()
			prodSnap := r.producer.Snapshot()

			rec := Record{
				Time:           now,
				Engine:         engSnap,
				Producer:       prodSnap,
				InstantQPS:     float64(prodSnap.SubmitTotal-lastProd.SubmitTotal) / r.interval.Seconds(),
				InstantDrop:    (engSnap.SubmitTotal - engSnap.SubmitSuccess) - (lastEng.SubmitTotal - lastEng.SubmitSuccess),
				InstantConsume: engSnap.ConsumeTotal - lastEng.ConsumeTotal,
			}

			r.mu.Lock()
			r.history = append(r.history, rec)
			r.mu.Unlock()

			r.renderRealtime(rec)

			lastEng = engSnap
			lastProd = prodSnap
		}
	}
}

func (r *Reporter) renderRealtime(rec Record) {
	dropTotal := rec.Engine.SubmitTotal - rec.Engine.SubmitSuccess
	fmt.Printf(
		"[realtime] %s | qps=%.0f success=%d drop=%d consume=%d queue=%d/%-6d latency=%.3fms\n",
		rec.Time.Format("15:04:05"),
		rec.InstantQPS,
		rec.Producer.SubmitSuccess,
		dropTotal,
		rec.Engine.ConsumeTotal,
		rec.Engine.QueueDepth,
		rec.Engine.QueueCapacity,
		rec.Engine.AvgLatencyMs,
	)
}

func (r *Reporter) renderFinalReport() {
	r.mu.Lock()
	records := make([]Record, len(r.history))
	copy(records, r.history)
	r.mu.Unlock()

	if len(records) == 0 {
		return
	}

	first := records[0]
	last := records[len(records)-1]

	totalDuration := last.Time.Sub(first.Time)
	totalProduced := last.Producer.SubmitTotal
	totalSuccess := last.Producer.SubmitSuccess
	totalDrop := last.Engine.SubmitDrop
	totalConsume := last.Engine.ConsumeTotal

	var maxQueue int64
	var totalLatency float64
	var dropPhases int

	for _, rec := range records {
		if rec.Engine.QueueDepth > maxQueue {
			maxQueue = rec.Engine.QueueDepth
		}
		totalLatency += rec.Engine.AvgLatencyMs
		if rec.InstantDrop > 0 {
			dropPhases++
		}
	}
	avgLatency := totalLatency / float64(len(records))

	successRate := float64(0)
	if totalProduced > 0 {
		successRate = float64(totalSuccess) / float64(totalProduced) * 100
	}

	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, "╔════════════════════════════════════════════════════════════════════╗")
	_, _ = fmt.Fprintln(os.Stderr, "║                      RUNTIME FINAL REPORT                          ║")
	_, _ = fmt.Fprintln(os.Stderr, "╠════════════════════════════════════════════════════════════════════╣")
	_, _ = fmt.Fprintf(os.Stderr, "║  Total Duration      : %-43s ║\n", totalDuration.Round(time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "║  Total Produced      : %-43d ║\n", totalProduced)
	_, _ = fmt.Fprintf(os.Stderr, "║  Total Success       : %-43d ║\n", totalSuccess)
	_, _ = fmt.Fprintf(os.Stderr, "║  Total Dropped       : %-43d ║\n", totalDrop)
	_, _ = fmt.Fprintf(os.Stderr, "║  Total Consumed      : %-43d ║\n", totalConsume)
	_, _ = fmt.Fprintf(os.Stderr, "║  Max Queue Depth     : %-43d ║\n", maxQueue)
	_, _ = fmt.Fprintf(os.Stderr, "║  Avg Consume Latency : %-43.3f ║\n", avgLatency)
	_, _ = fmt.Fprintf(os.Stderr, "║  Drop Phases         : %-43d ║\n", dropPhases)
	_, _ = fmt.Fprintf(os.Stderr, "║  Success Rate        : %-42.2f%% ║\n", successRate)
	_, _ = fmt.Fprintln(os.Stderr, "╚════════════════════════════════════════════════════════════════════╝")
}
