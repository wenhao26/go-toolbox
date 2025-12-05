package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 任务状态
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusCancelled  TaskStatus = "cancelled"
)

// 任务结构体
type Task struct {
	ID        string
	Data      interface{}
	Status    TaskStatus
	Result    interface{}
	Error     error
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 任务处理器接口
type TaskHandler interface {
	Process(ctx context.Context, task *Task) (interface{}, error)
}

// 任务管理器
type TaskManager struct {
	// 任务队列
	taskQueue chan *Task
	// 结果通道
	resultChan chan *Task
	// 取消信号通道
	cancelChan chan string
	// 任务映射
	tasks sync.Map
	// 工作池
	workerPool []*Worker
	// 配置
	maxWorkers int
	queueSize  int
	// 统计
	stats *Stats
	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// 工作协程
type Worker struct {
	id        int
	taskQueue chan *Task
	manager   *TaskManager
	handler   TaskHandler
}

// 统计信息
type Stats struct {
	Submitted     int64
	Completed     int64
	Failed        int64
	Cancelled     int64
	ActiveWorkers int64
	QueueLength   int64
}

// 模拟任务处理器
type SimulationHandler struct {
	minProcessingTime int
	maxProcessingTime int
	failureRate       float64
}

// 模拟任务结果
type SimulationResult struct {
	WorkerID       int
	ProcessingTime time.Duration
	DataSize       int
	Timestamp      time.Time
}

func main() {
	// 创建上下文，支持超时和取消
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建任务管理器
	manager := NewTaskManager(ctx, 100, 1000)

	// 启动任务管理器
	manager.Start()

	// 模拟提交任务
	go submitTasks(ctx, manager)

	// 监控任务执行情况
	go monitorTasks(ctx, manager)

	// 模拟取消一些任务
	go cancelSomeTasks(ctx, manager)

	// 等待上下文取消或超时
	<-ctx.Done()

	// 优雅关闭
	manager.Stop()

	// 打印最终统计
	manager.PrintStats()
}

// 创建新的任务管理器
func NewTaskManager(ctx context.Context, maxWorkers, queueSize int) *TaskManager {
	childCtx, cancel := context.WithCancel(ctx)

	// 创建模拟处理器（30%的失败率，处理时间100-500ms）
	_ = &SimulationHandler{
		minProcessingTime: 100,
		maxProcessingTime: 500,
		failureRate:       0.3,
	}

	return &TaskManager{
		taskQueue:  make(chan *Task, queueSize),
		resultChan: make(chan *Task, queueSize),
		cancelChan: make(chan string, 100),
		maxWorkers: maxWorkers,
		queueSize:  queueSize,
		stats:      &Stats{},
		ctx:        childCtx,
		cancel:     cancel,
		workerPool: make([]*Worker, 0, maxWorkers),
	}
}

// 启动任务管理器
func (tm *TaskManager) Start() {
	log.Println("启动任务管理器...")

	// 创建工作协程
	for i := 0; i < tm.maxWorkers; i++ {
		worker := NewWorker(i, tm.taskQueue, tm)
		tm.workerPool = append(tm.workerPool, worker)
		tm.wg.Add(1)
		go worker.Start(&tm.wg)
	}

	// 启动结果处理器
	tm.wg.Add(1)
	go tm.processResults(&tm.wg)

	// 启动取消处理器
	tm.wg.Add(1)
	go tm.processCancellations(&tm.wg)

	log.Printf("任务管理器已启动，工作协程数: %d", tm.maxWorkers)
}

// 停止任务管理器
func (tm *TaskManager) Stop() {
	log.Println("停止任务管理器...")

	// 发送取消信号
	tm.cancel()

	// 等待所有协程完成
	tm.wg.Wait()

	// 关闭通道
	close(tm.taskQueue)
	close(tm.resultChan)
	close(tm.cancelChan)

	log.Println("任务管理器已停止")
}

// 提交任务
func (tm *TaskManager) SubmitTask(data interface{}) (string, error) {
	// 生成任务ID
	taskID := fmt.Sprintf("task-%d-%d", time.Now().UnixNano(), rand.Intn(1000))

	task := &Task{
		ID:        taskID,
		Data:      data,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 存储任务
	tm.tasks.Store(taskID, task)

	// 更新统计
	atomic.AddInt64(&tm.stats.Submitted, 1)
	atomic.AddInt64(&tm.stats.QueueLength, 1)

	select {
	case tm.taskQueue <- task:
		// 成功提交到队列
		log.Printf("任务已提交: %s", taskID)
		return taskID, nil
	case <-tm.ctx.Done():
		// 上下文已取消
		task.Status = StatusCancelled
		task.UpdatedAt = time.Now()
		atomic.AddInt64(&tm.stats.Cancelled, 1)
		atomic.AddInt64(&tm.stats.QueueLength, -1)
		return "", tm.ctx.Err()
	default:
		// 队列已满
		task.Status = StatusFailed
		task.Error = fmt.Errorf("任务队列已满")
		task.UpdatedAt = time.Now()
		atomic.AddInt64(&tm.stats.Failed, 1)
		atomic.AddInt64(&tm.stats.QueueLength, -1)
		return "", task.Error
	}
}

// 取消任务
func (tm *TaskManager) CancelTask(taskID string) bool {
	select {
	case tm.cancelChan <- taskID:
		return true
	case <-tm.ctx.Done():
		return false
	default:
		log.Printf("取消队列已满，无法取消任务: %s", taskID)
		return false
	}
}

// 处理结果
func (tm *TaskManager) processResults(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task := <-tm.resultChan:
			// 更新任务状态
			tm.tasks.Store(task.ID, task)

			// 更新统计
			atomic.AddInt64(&tm.stats.QueueLength, -1)

			switch task.Status {
			case StatusCompleted:
				atomic.AddInt64(&tm.stats.Completed, 1)
				log.Printf("任务完成: %s, 结果: %v", task.ID, task.Result)
			case StatusFailed:
				atomic.AddInt64(&tm.stats.Failed, 1)
				log.Printf("任务失败: %s, 错误: %v", task.ID, task.Error)
			case StatusCancelled:
				atomic.AddInt64(&tm.stats.Cancelled, 1)
				log.Printf("任务取消: %s", task.ID)
			}

		case <-tm.ctx.Done():
			log.Println("结果处理器停止")
			return
		}
	}
}

// 处理取消请求
func (tm *TaskManager) processCancellations(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case taskID := <-tm.cancelChan:
			if val, ok := tm.tasks.Load(taskID); ok {
				task := val.(*Task)
				if task.Status == StatusPending || task.Status == StatusProcessing {
					task.Status = StatusCancelled
					task.UpdatedAt = time.Now()
					tm.resultChan <- task
					log.Printf("任务已取消: %s", taskID)
				}
			}

		case <-tm.ctx.Done():
			log.Println("取消处理器停止")
			return
		}
	}
}

// 获取统计信息
func (tm *TaskManager) GetStats() Stats {
	return Stats{
		Submitted:     atomic.LoadInt64(&tm.stats.Submitted),
		Completed:     atomic.LoadInt64(&tm.stats.Completed),
		Failed:        atomic.LoadInt64(&tm.stats.Failed),
		Cancelled:     atomic.LoadInt64(&tm.stats.Cancelled),
		ActiveWorkers: atomic.LoadInt64(&tm.stats.ActiveWorkers),
		QueueLength:   atomic.LoadInt64(&tm.stats.QueueLength),
	}
}

// 打印统计信息
func (tm *TaskManager) PrintStats() {
	stats := tm.GetStats()

	fmt.Println("\n=== 任务执行统计 ===")
	fmt.Printf("提交的任务数: %d\n", stats.Submitted)
	fmt.Printf("完成的任务数: %d\n", stats.Completed)
	fmt.Printf("失败的任务数: %d\n", stats.Failed)
	fmt.Printf("取消的任务数: %d\n", stats.Cancelled)
	fmt.Printf("当前队列长度: %d\n", stats.QueueLength)
	fmt.Printf("活跃工作协程: %d\n", stats.ActiveWorkers)

	totalProcessed := stats.Completed + stats.Failed + stats.Cancelled
	if stats.Submitted > 0 {
		successRate := float64(stats.Completed) / float64(stats.Submitted) * 100
		fmt.Printf("成功率: %.2f%%\n", successRate)
	}
	if totalProcessed > 0 {
		fmt.Printf("处理进度: %d/%d (%.1f%%)\n",
			totalProcessed, stats.Submitted,
			float64(totalProcessed)/float64(stats.Submitted)*100)
	}
}

// 创建新的工作协程
func NewWorker(id int, taskQueue chan *Task, manager *TaskManager) *Worker {
	return &Worker{
		id:        id,
		taskQueue: taskQueue,
		manager:   manager,
		handler: &SimulationHandler{
			minProcessingTime: 100,
			maxProcessingTime: 500,
			failureRate:       0.3,
		},
	}
}

// 启动工作协程
func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// 增加活跃工作协程计数
	atomic.AddInt64(&w.manager.stats.ActiveWorkers, 1)
	defer atomic.AddInt64(&w.manager.stats.ActiveWorkers, -1)

	log.Printf("工作协程 %d 启动", w.id)

	for {
		select {
		case task := <-w.taskQueue:
			// 处理任务
			w.processTask(task)

		case <-w.manager.ctx.Done():
			log.Printf("工作协程 %d 停止", w.id)
			return
		}
	}
}

// 处理单个任务
func (w *Worker) processTask(task *Task) {
	// 检查任务是否已被取消
	if task.Status == StatusCancelled {
		return
	}

	// 更新任务状态
	task.Status = StatusProcessing
	task.UpdatedAt = time.Now()

	log.Printf("工作协程 %d 开始处理任务: %s", w.id, task.ID)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(w.manager.ctx, 2*time.Second)
	defer cancel()

	// 处理任务
	result, err := w.handler.Process(ctx, task)

	// 更新任务状态
	if err != nil {
		task.Status = StatusFailed
		task.Error = err
	} else {
		task.Status = StatusCompleted
		task.Result = result
	}
	task.UpdatedAt = time.Now()

	// 发送到结果通道
	select {
	case w.manager.resultChan <- task:
		// 成功发送结果
	case <-w.manager.ctx.Done():
		// 上下文已取消
	}

	log.Printf("工作协程 %d 完成处理任务: %s", w.id, task.ID)
}

// 模拟任务处理器实现
func (sh *SimulationHandler) Process(ctx context.Context, task *Task) (interface{}, error) {
	// 模拟处理时间
	processingTime := rand.Intn(sh.maxProcessingTime-sh.minProcessingTime) + sh.minProcessingTime

	// 使用select监听上下文取消或超时
	select {
	case <-time.After(time.Duration(processingTime) * time.Millisecond):
		// 模拟处理完成

		// 模拟随机失败
		if rand.Float64() < sh.failureRate {
			return nil, fmt.Errorf("模拟处理失败")
		}

		// 返回模拟结果
		result := &SimulationResult{
			WorkerID:       rand.Intn(1000),
			ProcessingTime: time.Duration(processingTime) * time.Millisecond,
			DataSize:       rand.Intn(1024 * 1024), // 1MB以内
			Timestamp:      time.Now(),
		}

		return result, nil

	case <-ctx.Done():
		// 上下文被取消或超时
		return nil, fmt.Errorf("处理超时或被取消: %w", ctx.Err())
	}
}

// 模拟提交任务
func submitTasks(ctx context.Context, manager *TaskManager) {
	taskCounter := 0

	for {
		select {
		case <-time.After(50 * time.Millisecond):
			// 每50ms提交一个任务
			taskData := fmt.Sprintf("任务数据-%d", taskCounter)
			taskID, err := manager.SubmitTask(taskData)

			if err != nil {
				log.Printf("提交任务失败: %v", err)
			} else {
				taskCounter++
			}

			// 每提交10个任务，打印一次统计
			if taskCounter%10 == 0 {
				stats := manager.GetStats()
				log.Printf("任务ID %d, 已提交 %d 个任务, 队列长度: %d", taskID, taskCounter, stats.QueueLength)
			}

		case <-ctx.Done():
			log.Printf("任务提交器停止，共提交 %d 个任务", taskCounter)
			return
		}
	}
}

// 监控任务执行情况
func monitorTasks(ctx context.Context, manager *TaskManager) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := manager.GetStats()
			log.Printf("[监控] 统计: 提交=%d, 完成=%d, 失败=%d, 取消=%d, 队列=%d, 活跃=%d",
				stats.Submitted, stats.Completed, stats.Failed,
				stats.Cancelled, stats.QueueLength, stats.ActiveWorkers)

		case <-ctx.Done():
			log.Println("任务监控器停止")
			return
		}
	}
}

// 模拟取消一些任务
func cancelSomeTasks(ctx context.Context, manager *TaskManager) {
	time.Sleep(10 * time.Second) // 等待10秒后开始取消任务

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	cancelCount := 0

	for {
		select {
		case <-ticker.C:
			// 模拟取消最近提交的3个任务
			for i := 0; i < 3; i++ {
				taskID := fmt.Sprintf("task-%d", time.Now().UnixNano()-int64(rand.Intn(1000000000)))
				if manager.CancelTask(taskID) {
					cancelCount++
				}
			}
			log.Printf("已尝试取消 %d 个任务", cancelCount)

		case <-ctx.Done():
			log.Println("任务取消器停止")
			return
		}
	}
}
