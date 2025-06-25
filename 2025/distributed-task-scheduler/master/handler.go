package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux" // 使用 Gorilla Mux 作为路由处理器

	"toolbox/2025/distributed-task-scheduler/constant"
)

// - 实现任务的提交和结果查询的HTTP请求处理逻辑

// SubmitTaskHandler 处理任务提交请求
func SubmitTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

	// 解析请求体中的任务数据
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 生成任务ID
	task.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	task.Status = constant.Pending

	// 将任务添加到任务队列
	taskQueue <- &task

	// 返回任务提交成功的响应
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "task submitted",
		"task_id": task.ID,
	})
}

// GetResultHandler 处理任务结果查询请求
func GetResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"] // 从请求参数中获取任务ID

	mu.Lock()
	task, ok := taskMap[taskID] // 根据任务ID查找任务
	mu.Unlock()

	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// 返回任务的结果
	json.NewEncoder(w).Encode(task)
}
