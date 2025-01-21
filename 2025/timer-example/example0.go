// 模拟定时发布任务
package main

import (
	"fmt"
	"time"
)

// 任务逻辑
func runTask() {
	fmt.Println("Task is running at:", time.Now().Format("2006-01-02 15:04:05"))
	// TODO 这里实现具体的任务逻辑
}

func main() {
	inputTimeStr := "2025-01-21 14:44:00"
	// inputTime, err := time.Parse("2006-01-02 15:04:05", inputTimeStr)

	// 解析时间时指定本地时区
	loc, _ := time.LoadLocation("Local") // 加载本地时区
	inputTime, err := time.ParseInLocation("2006-01-02 15:04:05", inputTimeStr, loc)
	if err != nil {
		fmt.Println("Invalid time format. Please use 'YYYY-MM-DD HH:MM:SS'.")
		return
	}

	// 获取当前时间
	currentTime := time.Now()

	// 检查输入时间是否大于当前时间1分钟
	if inputTime.Before(currentTime.Add(1 * time.Minute)) {
		fmt.Println("The input time must be at least 1 minute later than the current time.")
		return
	}

	// 计算时间差
	duration := inputTime.Sub(currentTime)
	fmt.Printf("Time difference: %v\n", duration)

	// 启动定时器
	fmt.Printf("Task scheduled to run at: %s\n", inputTime.Format("2006-01-02 15:04:05"))
	timer := time.NewTimer(duration)

	// 等待定时器触发
	<-timer.C

	// 执行任务
	runTask()
}
