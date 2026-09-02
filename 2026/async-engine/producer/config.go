package producer

import (
	"time"
)

// Phase 定义一个流量阶段的持续时长和目标 QPS
type Phase struct {
	Duration time.Duration
	QPS      int
}

// Config 定义 Mock 生产者的运行参数
type Config struct {
	WorkerNum int     // 投递协程数量
	Phases    []Phase // 流量阶段序列，按顺序执行
}
