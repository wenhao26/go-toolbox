package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// - 支持 CRON 表达式的定时任务系统，可以更灵活地定义任务的执行时间

// Task 定义任务结构
type Task struct {
	Name  string
	Func  func()
	Entry cron.EntryID
	Cron  *cron.Cron
}

// TaskManager 定义任务管理器
type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]*Task
}

// NewTaskManager 创建任务管理器实例
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// AddTask 添加任务
func (m *TaskManager) AddTask(name, spec string, taskFunc func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tasks[name]; ok {
		return fmt.Errorf("task %s already exists", name)
	}

	// 创建一个新的cron实例
	c := cron.New()

	// 添加任务到cron
	entryID, err := c.AddFunc(spec, taskFunc)
	if err != nil {
		return err
	}

	// 启动cron
	c.Start()

	task := &Task{
		Name:  name,
		Func:  taskFunc,
		Entry: entryID,
		Cron:  c,
	}

	m.tasks[name] = task
	fmt.Printf("task %s added with spec:%s\n", name, spec)

	return nil
}

// StopTask 停止任务
func (m *TaskManager) StopTask(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[name]
	if !ok {
		fmt.Printf("task %s does not exist\n", name)
	}

	// 停止任务
	task.Cron.Remove(task.Entry)

	// 停止cron实例
	task.Cron.Stop()

	delete(m.tasks, name)
	fmt.Printf("task %s stopped\n", name)
}

// GetTasks 获取所有任务
func (m *TaskManager) GetTasks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("Current tasks:")
	for name := range m.tasks {
		fmt.Println(name)
	}
}

func main() {
	manager := NewTaskManager()

	// 添加任务
	var err error

	err = manager.AddTask("task01", "*/3 * * * *", func() {
		fmt.Println("task01 is running at ", time.Now())
	})
	if err != nil {
		fmt.Println("error adding task01:", err)
		return
	}

	err = manager.AddTask("task02", "*/5 * * * *", func() {
		fmt.Println("task02 is running at ", time.Now())
	})
	if err != nil {
		fmt.Println("error adding task02:", err)
		return
	}

	// 列出任务
	manager.GetTasks()

	// 模拟运行一分钟后停止任务
	time.Sleep(30 * time.Second)
	manager.StopTask("task01")

	manager.GetTasks()
	time.Sleep(1 * time.Minute)
}
