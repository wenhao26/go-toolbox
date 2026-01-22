package handlers

import (
	"fmt"
	"net/http"
	"toolbox/2026/cron-jobs/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务操作者结构体
type TaskHandler struct {
	Mgr *service.CronManager
}

// RequestParams 定义请求参数
type RequestParams struct {
	Name string `json:"name" binding:"required"`
	Spec string `json:"spec"`
}

// RegisterRoutes 注册API路由
func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/v1/tasks")
	{
		group.POST("/add", h.AddTask)   // 启动任务
		group.POST("/stop", h.StopTask) // 停止任务
		group.GET("/list", h.ListTasks) // 查看运行中任务
	}
}

// AddTask 启动任务
func (h *TaskHandler) AddTask(c *gin.Context) {
	var req RequestParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 模拟具体的业务逻辑函数
	job := func() {
		fmt.Printf("[Job Executed] 任务名: %s\n", req.Name)
	}

	if err := h.Mgr.AddTask(req.Name, req.Spec, job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务启动成功"})
}

// StopTask 停止任务
func (h *TaskHandler) StopTask(c *gin.Context) {
	var req RequestParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Mgr.RemoveTask(req.Name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务已停止"})
}

// ListTasks 查看运行中任务
func (h *TaskHandler) ListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": h.Mgr.ListTasks()})
}
