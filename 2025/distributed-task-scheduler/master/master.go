package master

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"toolbox/2025/distributed-task-scheduler/constant"
)

// - 实现主节点的启动和停止逻辑

// 定义任务队列和任务映射表
var (
	taskQueue = make(chan *Task, 100)  // 任务队列，用于存储待处理的任务
	taskMap   = make(map[string]*Task) // 任务映射表，用于根据任务ID查找任务
	mu        sync.Mutex               // 互斥锁，用于保护taskMap的并发访问
)

// StartMaster 启动主节点
func StartMaster() error {
	// 启动协程来处理任务队列
	go handleTaskQueue()

	// 注册HTTP请求处理函数
	http.HandleFunc("/submit", SubmitTaskHandler)
	http.HandleFunc("/result", GetResultHandler)

	// 启动HTTP服务器
	go func() {
		log.Println("master node started on:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("failed to start master node:%v", err)
		}
	}()

	return nil
}

// StopMaster 停止主节点
func StopMaster() {
	close(taskQueue) // 关闭任务队列
}

// handleTaskQueue 处理任务队列中的任务
func handleTaskQueue() {
	for task := range taskQueue {
		mu.Lock()
		taskMap[task.ID] = task // 将任务添加到任务映射表
		mu.Unlock()

		// 模拟任务处理过程
		time.Sleep(5 * time.Second)
		task.Result = fmt.Sprintf("result of %s", task.Content)
		task.Status = constant.Completed

		mu.Lock()
		taskMap[task.ID] = task // 更新任务的状态和结果
		mu.Unlock()
	}
}
