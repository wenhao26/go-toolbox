package engine

import (
	"fmt"
	"math/rand"
	"time"
	"toolbox/2026/async-engine/model"
)

type worker struct {
	id     int
	queue  <-chan *model.Task
	engine *Engine
}

func (w *worker) run() {
	defer w.engine.wg.Done()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(w.id)))

	for task := range w.queue {
		w.process(task, rnd)
	}
}

func (w *worker) process(task *model.Task, rnd *rand.Rand) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[worker-%d] panic recovered: %v\n", w.id, r)
		}
	}()

	start := time.Now()

	// Mock: 日志序列化与落盘前的处理，平均耗时约 200 微秒
	_ = fmt.Sprintf(
		"trace=%s method=%s path=%s status=%d latency=%s ip=%s size=%d",
		task.TraceID,
		task.Method,
		task.Path,
		task.Status,
		task.Latency,
		task.ClientIP,
		task.PayloadSize,
	)

	// Mock: IO 波动 -> [150 ~ 250] 微秒
	time.Sleep(time.Duration(150+rnd.Intn(100)) * time.Microsecond)

	w.engine.recordConsume(time.Since(start))
}
