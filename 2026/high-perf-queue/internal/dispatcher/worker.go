package dispatcher

import (
	"log"
	"runtime/debug"
	"sync"
)

// Worker 代表一个独立的消费者协程
type Worker struct {
	id        int             // Worker 唯一标识符
	taskQueue chan *Task      // 共享的任务通道
	wg        *sync.WaitGroup // 用于通知调度器自身已安全退出
}

// NewWorker 创建一个消费者实例
func NewWorker(id int, taskQueue chan *Task, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:        id,
		taskQueue: taskQueue,
		wg:        wg,
	}
}

// Start 启动消费者
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		// 严密捕获消费者自身或业务代码触发的 panic，保障宿主进程不崩溃
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Worker-%d] 严重错误: 捕获到未处理的 Panic: %v\n堆栈信息:\n%s", w.id, err, string(debug.Stack()))
			}
		}()

		// 阻塞监听通道，直到通道被 close 且数据被拉取完毕
		for task := range w.taskQueue {
			w.execute(task)
		}
		log.Printf("[Worker-%d] 任务队列已空且通道已关闭，安全退出。\n", w.id)
	}()
}

// execute 处理单条任务
func (w *Worker) execute(task *Task) {
	if task == nil {
		return
	}

	// 无论业务执行成功、失败或发生 panic，都必须确保 Task 被回收
	defer ReleaseTask(task)

	// 再次做一层业务级的 panic recovery，防止单条任务挂掉影响当前 Worker 的后续消费
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[Worker-%d] 任务执行失败 [Payload: %v]: 触发 Panic: %v", w.id, task.Payload, err)
		}
	}()

	// 检查上游上下文是否已经超时或取消
	if err := task.Ctx.Err(); err != nil {
		log.Printf("[Worker-%d] 任务跳过 [Payload: %v]: 上下文已取消: %v\n", w.id, task.Payload, err)
		return
	}

	// 执行标准业务 Handler
	if err := task.Handler(task.Ctx, task.Payload); err != nil {
		// 生产环境通常在这里进行重试逻辑（Retry）或投递到死信队列（DLQ）
		log.Printf("[Worker-%d] 业务处理返回错误 [Payload: %v]: %v\n", w.id, task.Payload, err)
	}

	log.Printf("[Worker-%d] 任务消费成功 [Payload: %v]", w.id, task.Payload)
}
