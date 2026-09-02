package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toolbox/2026/async-engine/engine"
	"toolbox/2026/async-engine/producer"
	"toolbox/2026/async-engine/stats"
)

func main() {
	// 引擎配置：6 个消费协程，队列容量 50000。
	// 单协程平均处理耗时约 200 微秒，理论吞吐约 30000/s，刚好覆盖目标峰值。
	engCfg := engine.Config{
		WorkerNum: 20,
		QueueSize: 50000,
	}
	eng := engine.NewEngine(engCfg)
	eng.Start()

	// 生产者配置：100 个投递协程，分阶段动态调整流量。
	// 阶段设计覆盖稳态、目标峰值、恢复、超限压测与持续峰值场景。
	prodCfg := producer.Config{
		WorkerNum: 100,
		Phases: []producer.Phase{
			{Duration: 5 * time.Minute, QPS: 10000},  // 稳态低峰，观察基线
			{Duration: 3 * time.Minute, QPS: 30000},  // 目标峰值，系统饱和边缘
			{Duration: 5 * time.Minute, QPS: 10000},  // 恢复阶段，队列消化
			{Duration: 2 * time.Minute, QPS: 50000},  // 超限压测，触发背压丢弃
			{Duration: 15 * time.Minute, QPS: 30000}, // 持续峰值，长期稳定性
		},
	}
	prod := producer.NewProducer(prodCfg, eng)
	prod.Start()

	// 报表协程：每 5 秒采样一次，输出实时看板。
	reporter := stats.NewReporter(eng, prod, 5*time.Second)
	reporter.Start()

	fmt.Println("[main] system started, press Ctrl+C or wait 30min to stop")

	// 监听中断信号，同时设置 30 分钟自动停止。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	autoStop := time.NewTimer(30 * time.Minute)
	defer autoStop.Stop()

	select {
	case <-sigCh:
		fmt.Println("[main] received shutdown signal")
	case <-autoStop.C:
		fmt.Println("[main] reached 30min auto stop")
	}

	// 优雅退出顺序：先停生产（不再产生新任务），再停消费（消化剩余队列），最后输出报表。
	prod.Stop()
	eng.Stop()
	reporter.Stop()

	fmt.Println("[main] system shutdown complete")
}
