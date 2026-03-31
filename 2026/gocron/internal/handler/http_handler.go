package handler

import (
	"fmt"
	"net/http"
	"toolbox/2026/gocron/internal/model"
	"toolbox/2026/gocron/internal/service"

	"github.com/gin-gonic/gin"
)

// -- Web 接口处理 --

// SaveTask 保存并启动任务
func SaveTask(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("参数错误：%v", err.Error()),
		})
		return
	}

	// 写入持久化存储 Redis
	err := service.SaveTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("保存任务到Redis错误：%v", err.Error()),
		})
		return
	}

	// 更新内存中的调度器
	err = service.GlobalTaskScheduler.AddOrUpdate(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("更新调度器错误：%v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "ID不能为空",
		})
		return
	}

	// 从调度引擎停止任务
	service.GlobalTaskScheduler.Remove(id)

	// 从 Redis 彻底删除
	service.DeleteTask(id)

	c.JSON(200, gin.H{"status": "ok"})
}

// GetTaskList 获取任务列表
func GetTaskList(c *gin.Context) {
	tasks, _ := service.LoadAllTasks()
	c.JSON(http.StatusOK, tasks)
}

// GetLogs 获取任务执行日志
func GetLogs(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, service.GetLogs(id))
}
