package scheduler

import (
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
)

// Task 定义定时任务
type Task struct {
	ID     cron.EntryID // CRON返回的任务ID
	Spec   string       // CRON表达式
	Func   func()       // 要执行的任务函数
	Active bool         // 当前任务是否处于激活状态
}

// Scheduler 调度器
type Scheduler struct {
	cron  *cron.Cron       // CRON实例
	tasks map[string]*Task // 已注册任务，按名称索引
	mu    sync.Mutex       // 互斥锁，确保并发安全
}

// NewScheduler 创建一个支持秒级的调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:  cron.New(cron.WithSeconds()), // 启用秒级调度
		tasks: make(map[string]*Task),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(name, spec string, fn func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[name]; exists {
		return errors.New("任务已经存在")
	}

	// 添加任务CRON调度器中
	entryID, err := s.cron.AddFunc(spec, fn)
	if err != nil {
		return err
	}

	// 保存任务信息
	s.tasks[name] = &Task{
		ID:     entryID,
		Spec:   spec,
		Func:   fn,
		Active: true,
	}

	return nil
}

// StopTask 停止（移除）某个正在运行的任务
func (s *Scheduler) StopTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[name]
	if !ok || !task.Active {
		return errors.New("任务不存在或已停止")
	}

	s.cron.Remove(task.ID)
	task.Active = false

	return nil
}

// RestartTask 重启一个已停止的任务（再次添加到调度器）
func (s *Scheduler) RestartTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[name]
	if !ok {
		return errors.New("任务不存在，重启失败")
	}
	if task.Active {
		return errors.New("任务正在运行中")
	}

	entryID, err := s.cron.AddFunc(task.Spec, task.Func)
	if err != nil {
		return err
	}

	task.ID = entryID
	task.Active = true

	return nil
}

// RemoveTask 永久删除任务（从调度器和任务列表中移除）
func (s *Scheduler) RemoveTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[name]
	if !ok {
		return errors.New("任务不存在，删除失败")
	}
	if task.Active {
		s.cron.Remove(task.ID)
	}

	delete(s.tasks, name)

	return nil
}

// ListTasks 返回所有任务的状态（true：运行中，false：已停止）
func (s *Scheduler) ListTasks() map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	status := make(map[string]bool)
	for name, task := range s.tasks {
		status[name] = task.Active
	}

	return status
}
