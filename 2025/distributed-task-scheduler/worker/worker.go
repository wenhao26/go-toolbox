package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"toolbox/2025/distributed-task-scheduler/constant"
)

// - 实现工作节点的启动和停止逻辑

// StartWorker 启动工作节点
func StartWorker() error {
	// 注册任务处理请求的处理函数
	http.HandleFunc("/process", processTaskHandler)

	// 启动HTTP服务器
	go func() {
		log.Println("worker node started on:8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("failed to start worker node: %v", err)
		}
	}()

	return nil
}

// StopWorker 停止工作节点
func StopWorker() {
	// TODO 清理资源（如果有需要）
}

// processTaskHandler 处理任务请求
func processTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

	// 解析请求体中的任务数据
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 模拟任务处理过程
	time.Sleep(5 * time.Second)
	task.Result = fmt.Sprintf("result of %s", task.Content)
	task.Status = constant.Completed

	// 返回任务的处理结果
	json.NewEncoder(w).Encode(task)
}
