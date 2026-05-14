package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Task 任务信息
type Task struct {
	ID      int
	Payload string
}

// Result 处理结果
type Result struct {
	TaskID int
	Data   string
	Err    error
}

// Engine 处理引擎
type Engine struct {
	taskChan   chan Task      // 任务分发通道
	resultChan chan Result    // 结果收集通道
	workerNum  int            // 开启工作数量（协程数）设置
	wg         sync.WaitGroup // 并发控制
}

// NewEngine 创建一个新的处理引擎实例
func NewEngine(workerNum int) *Engine {
	return &Engine{
		taskChan:   make(chan Task, 100), // 带缓冲，防止瞬间高并发卡死
		resultChan: make(chan Result, 100),
		workerNum:  workerNum,
	}
}

// Start 启动执行引擎
func (e *Engine) Start(ctx context.Context) {
	// 在启动处理引擎时，预先开启固定数量的协程
	for i := 0; i < e.workerNum; i++ {
		e.wg.Add(1)
		go e.worker(i, ctx)
	}

	// 开协程处理结果，不影响工作者逻辑
	go e.resultProcessor()
}

// worker 业务处理工作者
func (e *Engine) worker(id int, ctx context.Context) {
	defer e.wg.Done()

	fmt.Printf("[worker %d] 启动，开始待命...\n]", id)

	for {
		select {
		case task, ok := <-e.taskChan:
			if !ok {
				return
			}

			// 执行任务并带上“超时控制”
			res := e.handleTaskWithTimeout(task)
			e.resultChan <- res
		case <-ctx.Done():
			fmt.Printf("[worker %d] 收到退出信号，停止处理。\n", id)
			return
		}
	}
}

// handleTaskWithTimeout 处理任务（展示 select 经典用法 - 超时控制）
func (e *Engine) handleTaskWithTimeout(task Task) Result {
	// TODO：模拟一个带超时的处理过程
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 用一个临时通道获取结果
	done := make(chan Result, 1)

	// 模拟耗时业务逻辑
	go func() {
		sleepTime := time.Duration(rand.Intn(3000)) * time.Millisecond
		time.Sleep(sleepTime)

		done <- Result{TaskID: task.ID, Data: fmt.Sprintf("处理结果: %s", task.Payload)}
	}()

	// 使用 select 监听（核心）：要么成功，要么超时
	select {
	case res := <-done:
		return res
	case <-ctx.Done():
		return Result{TaskID: task.ID, Err: errors.New("任务执行超时")}
	}
}

// resultProcessor 持续处理结果集（消费者）
func (e *Engine) resultProcessor() {
	for result := range e.resultChan {
		if result.Err != nil {
			fmt.Printf("❌ 任务 %d 失败: %v\n", result.TaskID, result.Err)
		} else {
			fmt.Printf("✅ 任务 %d 成功: %s\n", result.TaskID, result.Data)
		}
	}
}

// Submit 提交任务
func (e *Engine) Submit(task Task) {
	e.taskChan <- task
}

func (e *Engine) Close() {
	close(e.taskChan) // 关闭任务通道
	e.wg.Wait()       // 等到所有 worker 处理完
	close(e.resultChan)

	fmt.Println("引擎已完全关闭。")
}

func main() {
	engine := NewEngine(50) // 开启100个并发处理协程

	// 控制生命周期 Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine.Start(ctx)

	// 模拟外部请求
	go func() {
		for i := 0; i < 2000; i++ {
			engine.Submit(Task{
				ID:      i,
				Payload: fmt.Sprintf("数据-%d", i),
			})
			time.Sleep(time.Millisecond * 50)
		}
	}()

	// 优雅退出：监听系统强制中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // 阻塞，直到按下【ctrl+c】

	fmt.Println("\n正在关闭系统...")
	engine.Close() // 安全清理资源
}
