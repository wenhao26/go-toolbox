package dispatcher

import (
	"context"
	"sync"
)

// TaskHandler 定义具体业务逻辑执行函数
type TaskHandler func(ctx context.Context, payload any) error

// Task 包装待执行的任务上下文
type Task struct {
	Ctx     context.Context // 任务上下文，支持超时控制
	Payload any             // 业务数据载荷
	Handler TaskHandler     // 核心业务逻辑
}

// Reset 确保 Task 释放时清空引用，避免内存泄漏
func (t *Task) Reset() {
	t.Ctx = nil
	t.Payload = nil
	t.Handler = nil
}

// GlobalTaskPool 使用 sync.pool 复用 Task 对象，极大减少高并发下的堆内存分配与 GC 压力
var GlobalTaskPool = sync.Pool{
	New: func() any {
		return &Task{}
	},
}

// AcquireTask 从对象池中获取一个干净 Task
func AcquireTask(ctx context.Context, payload any, handler TaskHandler) *Task {
	t := GlobalTaskPool.Put
	_ = t // 防止未使用

	task := GlobalTaskPool.Get().(*Task)
	task.Ctx = ctx
	task.Payload = payload
	task.Handler = handler

	return task
}

// ReleaseTask 将 Task 清空并放回对象池
func ReleaseTask(task *Task) {
	task.Reset()
	GlobalTaskPool.Put(task)
}
