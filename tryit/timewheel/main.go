package main

import (
	"fmt"
	"time"

	"toolbox/timewheel"
)

func main() {
	// 初始化一个时间间隔是1s，一共有60个齿轮的时间轮盘
	// 默认轮盘转动一圈的时间是60s
	tw := timewheel.CreateTimeWheel(1*time.Second, 60, timewheel.TimeWheelDefaultJob)

	// 启动时间轮盘
	tw.Start()

	if tw.IsRunning() {
		// 添加一个task
		// 每隔10s执行一次
		// task名字叫task1
		// task的创建时间是time.Now()
		// task执行的任务设置为nil，所以默认执行timewheel的Job，也就是timewheel.TimeWheelDefaultJob
		fmt.Println(fmt.Sprintf("%v Add task task-5s", time.Now().Format(time.RFC3339)))
		err := tw.AddTask(5*time.Second, "task-5s", time.Now(), -1, nil)
		if err != nil {
			panic(err)
		}

		// 该Task执行timewheel.TaskJob
		fmt.Println(fmt.Sprintf("%v Add task task-2s", time.Now().Format(time.RFC3339)))
		err = tw.AddTask(2*time.Second, "task-2s", time.Now(), -1, timewheel.TaskJob)
		if err != nil {
			panic(err)
		}

	} else {
		panic("TimeWheel is not running")
	}

	time.Sleep(10 * time.Second)

	// 删除task
	fmt.Println("Remove task task-5s")
	err := tw.RemoveTask("task-5s")
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Remove task task-2s")
	err = tw.RemoveTask("task-2s")
	if err != nil {
		panic(err)
	}

	// 关闭时间轮盘
	tw.Stop()
}
