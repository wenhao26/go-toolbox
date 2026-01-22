package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

// CronManager 定时任务管理器结构体
type CronManager struct {
	scheduler *cron.Cron
	entryIDs  sync.Map // key: 任务名称，value: cron.EntryID
}

// NewCronManager 创建一个新的定时任务管理器
func NewCronManager() *CronManager {
	c := cron.New(cron.WithSeconds())
	c.Start()

	return &CronManager{
		scheduler: c,
	}
}

// AddTask 添加并启动任务
func (m *CronManager) AddTask(name, spec string, cmd func()) error {
	if _, ok := m.entryIDs.Load(name); ok {
		return errors.New("任务已存在")
	}

	id, err := m.scheduler.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("cron表达式错误: %v", err)
	}

	m.entryIDs.Store(name, id)

	return nil
}

// RemoveTask 停止并删除任务
func (m *CronManager) RemoveTask(name string) error {
	v, ok := m.entryIDs.Load(name)
	if !ok {
		return errors.New("任务不存在")
	}

	id := v.(cron.EntryID)
	m.scheduler.Remove(id)
	m.entryIDs.Delete(name)

	return nil
}

// ListTasks 获取当前所有运行中的任务名称
func (m *CronManager) ListTasks() []string {
	var names []string
	m.entryIDs.Range(func(key, value any) bool {
		names = append(names, key.(string))
		return true
	})

	return names
}
