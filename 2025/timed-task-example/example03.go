package main

import (
	"fmt"
	"time"

	"toolbox/2025/timed-task-example/scheduler"
)

// - 一个生产可用级别的 Go 定时任务调度器组件
// - 支持功能：
// 	- 秒级调度（基于 robfig/cron/v3）
//	- 添加任务（命名任务，便于管理）
//	- 停止任务
//	- 重启任务
//	- 删除任务
//	- 查看所有任务状态（可扩展）

func main() {
	s := scheduler.NewScheduler()

	s.Start()
	_ = s.AddTask("print-time", "* * * * * *", func() {
		fmt.Println("执行时间：", time.Now())
	})

	// 模拟暂停任务
	time.Sleep(5 * time.Second)
	_ = s.StopTask("print-time")
	fmt.Println("print-time 已暂停")

	// 模拟重启任务
	time.Sleep(5 * time.Second)
	_ = s.RestartTask("print-time")
	fmt.Println("print-time 已重启")

	// 模拟删除任务
	time.Sleep(10 * time.Second)
	_ = s.RemoveTask("print-time")
	fmt.Println("print-time 已删除")

	select {}
}
