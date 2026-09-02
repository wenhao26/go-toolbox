package model

import (
	"time"
)

// Task 模拟 API 访问产生的日志任务
type Task struct {
	TraceID     string
	Timestamp   time.Time
	Method      string
	Path        string
	Status      int
	Latency     time.Duration
	ClientIP    string
	PayloadSize int
}
