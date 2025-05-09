package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task 任务结构体
type Task struct {
	fn    func() error
	retry int
}

func NewTask(fn func() error) *Task {
	return &Task{fn: fn, retry: 3}
}

// Scheduler 调度器结构体
type Scheduler struct {
	tasks       chan *Task
	wg          sync.WaitGroup
	rateLimiter <-chan time.Time
}

func NewScheduler(concurrency int, ratePerSecond int) *Scheduler {
	s := &Scheduler{
		tasks:       make(chan *Task, 100),
		rateLimiter: time.Tick(time.Second / time.Duration(ratePerSecond)),
	}

	for i := 0; i < concurrency; i++ {
		go s.worker()
	}

	return s
}

func (s *Scheduler) Submit(task *Task) {
	s.wg.Add(1)
	s.tasks <- task
}

func (s *Scheduler) worker() {
	for task := range s.tasks {
		<-s.rateLimiter // 限速

		err := task.fn()
		if err != nil && task.retry > 0 {
			fmt.Println("任务失败，重试中...")
			task.retry--
			s.Submit(task)
		} else if err != nil {
			fmt.Println("任务最终失败:", err)
		}

		s.wg.Done()
	}
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
	close(s.tasks)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 初始化调度器：2个并发worker，限速5次每秒
	scheduler := NewScheduler(2, 5)

	// 提交多个任务
	for i := 1; i <= 10; i++ {
		i := i // 避免闭包坑
		task := NewTask(func() error {
			fmt.Printf("执行任务 #%d\n", i)
			// 模拟失败概率
			if rand.Intn(100) < 30 {
				return errors.New("模拟失败")
			}
			fmt.Printf("任务 #%d 成功 ✅\n", i)
			return nil
		})
		scheduler.Submit(task)
	}

	// 等待所有任务完成
	scheduler.Wait()
	fmt.Println("所有任务完成 ✅")
}
