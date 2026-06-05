/*
== 高性能 AI 创作任务调度引擎 ==

本模块实现了基于 Go 并发三剑客（Goroutine, Channel, Select）的生产级任务分发网关。

核心功能点：

 1. 绝对优先级调度 (Absolute Priority Scheduling):
    利用 select 嵌套预检机制，强制 Worker 优先处理加急通道（highChan），
    有效解决批量生成任务阻塞用户即时预览请求的痛点。

 2. 并发限流与资源保护 (Rate Limiting & Safety):
    通过信号量（Buffered Channel）限制全局 AI 接口调用频率，防止触发服务商 RPM 限制，
    并利用 context.WithTimeout 确保链路级超时止损。

 3. 优雅停机方案 (Graceful Shutdown):
    支持 stopChan 广播式信号退出。收到关闭指令后，Worker 会确保当前领取的
    AI 任务完整执行并入库后方才退出，避免 Token 资源浪费。

使用建议：
- 推荐根据 AI 服务商的并发限制设置 NewProcessor 的 maxConcurrent 参数。
- 加急任务通过 Submit(Task{IsUrgent: true}) 提交，普通任务则设为 false。

设计哲学：
“不要通过共享内存来通信，而要通过通信来共享内存。”
本引擎追求对任务流向的绝对掌控感，而非简单的并发执行。
*/
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Task 业务任务结构
type Task struct {
	ID       int
	Title    string
	IsUrgent bool // 是否属于紧急任务
}

// Processor 核心调度器
type Processor struct {
	highChan chan Task     // 高优先级通道
	lowChan  chan Task     // 低优先级通道
	stopChan chan struct{} // 停止信号
	limit    chan struct{} // 最大并发限制

	wg sync.WaitGroup
}

func NewProcessor(maxConcurrent int) *Processor {
	return &Processor{
		highChan: make(chan Task, 100),
		lowChan:  make(chan Task, 1000),
		stopChan: make(chan struct{}),
		limit:    make(chan struct{}, maxConcurrent),
	}
}

// Start 启动常驻 Worker
func (p *Processor) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// Stop 优雅关闭
func (p *Processor) Stop() {
	close(p.stopChan) // 广播关闭信号
	p.wg.Wait()
	log.Println("✅ [System] 所有 Worker 已安全退出，引擎关闭")
}

// Submit 提交任务
func (p *Processor) Submit(task Task) {
	if task.IsUrgent {
		p.highChan <- task
	} else {
		p.lowChan <- task
	}
}

// worker 内部调度工作逻辑
func (p *Processor) worker(workerID int) {
	defer p.wg.Done()
	log.Printf("👷 Worker %d 已启动...", workerID)

	for {
		// 强制优先级检查
		// select 默认是随机的，为了实现“绝对优先级”，必须尝试只读 highChan
		select {
		case <-p.stopChan:
			return
		case task := <-p.highChan:
			p.handle(workerID, task)
			return
		default:
			// 如果执行到这里，说明当前瞬间没有高优先级任务处理
		}

		// 二次多路复用
		// 此处才考虑低优先级任务，但同时也要监听 highChan 和 stopChan
		select {
		case <-p.stopChan:
			return
		case task := <-p.highChan:
			p.handle(workerID, task)
		case task := <-p.lowChan:
			p.handle(workerID, task)
		case <-time.After(10 * time.Second):
			log.Printf("Worker %d 正在待命...", workerID)
		}
	}
}

// handle 处理具体业务逻辑
func (p *Processor) handle(workerID int, task Task) {
	p.limit <- struct{}{}        // 占用并发位
	defer func() { <-p.limit }() // 释放并发位

	timeout := 20 * time.Second
	if task.IsUrgent { // 紧急任务，响应更快处理
		timeout = 5 * time.Second
	}

	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	log.Printf("🛠️ [Worker %d] 开始处理任务: %s (加急: %v)", workerID, task.Title, task.IsUrgent)

	// 模拟 AI 耗时请求
	err := p.mockAICall(ctx, task)
	if err != nil {
		log.Printf("❌ [Worker %d] 任务失败: %s, 错误: %v", workerID, task.Title, err)
		return
	}
	log.Printf("✨ [Worker %d] 任务完成: %s", workerID, task.Title)
}

// mockAICall 模拟 AI 耗时请求
func (p *Processor) mockAICall(ctx context.Context, task Task) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	engine := NewProcessor(3)
	engine.Start(2)

	// 模拟提交任务
	go func() {
		// 优先级低任务
		for i := 0; i < 5; i++ {
			engine.Submit(Task{
				ID:       i,
				Title:    fmt.Sprintf("批量处理小说章节 %d", i),
				IsUrgent: false,
			})
		}

		// 模拟延迟几秒后，突然插入高优先级任务
		time.Sleep(3 * time.Second)
		engine.Submit(Task{
			ID:       99999,
			Title:    "加急任务-用户预览章节",
			IsUrgent: true,
		})
	}()

	time.Sleep(1 * time.Minute)
	engine.Stop()
}
