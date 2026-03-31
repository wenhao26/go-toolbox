package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
	"toolbox/2026/gocron/internal/model"

	"github.com/robfig/cron/v3"
)

// -- 调度引擎逻辑 --

// TaskScheduler 任务调度器
type TaskScheduler struct {
	cron  *cron.Cron
	nodes sync.Map // 存储 TaskID -> cron.EntryID
}

var GlobalTaskScheduler *TaskScheduler // 全局任务调度器

// InitScheduler 初始化调度器
func InitScheduler() {
	c := cron.New(cron.WithSeconds()) // 开启秒级解析
	GlobalTaskScheduler = &TaskScheduler{
		cron: c,
	}
	c.Start()
}

// AddOrUpdate 添加/更新调度任务
func (s *TaskScheduler) AddOrUpdate(task model.Task) error {
	s.Remove(task.ID) // 先清理旧的任务ID

	// 如果任务状态为停止，则不添加新的调度
	if task.Status == 0 {
		return nil
	}

	entryID, err := s.cron.AddFunc(task.Expr, func() {
		// execute(t)

		// 开启协程，确保不阻塞 cron 内部的调度循环
		//
		// 将“时间轮询”与“任务执行”彻底分离
		// 意味着无论任务 execute 耗时多久（甚至是阻塞操作），都不会拖慢系统对下一个任务的精准触发
		go func(t model.Task) {
			// 核心防御：捕获该协程可能出现的 panic，防止主进程崩溃
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Scheduler] Task %s panic recovered: %v", t.ID, r)
				}
			}()

			execute(t)
		}(task)
	})
	if err != nil {
		log.Printf("[Scheduler] AddFunc Error for task %s: %v", task.ID, err)
		return err
	}

	// 将新的调度 EntryID 存入同步 Map
	s.nodes.Store(task.ID, entryID)
	return nil
}

// Remove 移除任务
func (s *TaskScheduler) Remove(id string) {
	if val, ok := s.nodes.Load(id); ok {
		s.cron.Remove(val.(cron.EntryID))
		s.nodes.Delete(id)
	}
}

// execute 执行具体的Shell脚本
func execute(task model.Task) {
	start := time.Now()

	// 设置超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.Timeout)*time.Second)
	defer cancel()

	var buf bytes.Buffer

	cmd := exec.CommandContext(ctx, "bash", "-c", task.Command)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	// 记录日志到 Redis
	logEntry := model.Log{
		StartTime: start.Format("2006-01-02 15:04:05"),
		EndTime:   time.Now().Format("2006-01-02 15:04:05"),
		Output:    buf.String(),
		Success:   err == nil,
	}
	if err != nil {
		logEntry.Output += fmt.Sprintf("\n[Error]: %v", err)
	}

	// 保存日志
	SaveLog(task.ID, logEntry)
}
