package handler

import (
	"context"
	"encoding/json"
	"toolbox/2026/novel-creation/pkg/logger"

	"go.uber.org/zap"
)

// Task 定义通用的任务结构，可根据实际业务扩展
type Task struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Handle 处理单个任务
func Handle(ctx context.Context, taskData string) error {
	var task Task
	if err := json.Unmarshal([]byte(taskData), &task); err != nil {
		logger.Log.Error("invalid task json", zap.String("task", taskData), zap.Error(err))
		// 格式错误的任务无需重试，直接返回 nil 表示“已消费但丢弃”
		return nil
	}

	logger.Log.Info("processing task", zap.String("task_id", task.ID), zap.String("type", task.Type))

	// 根据任务类型分发到不同业务逻辑
	switch task.Type {
	case "story":
		return handleStoryTask(ctx, task)
	case "poem":
		return handlePoemTask(ctx, task)
	default:
		logger.Log.Warn("unknown task type", zap.String("type", task.Type))
		return nil // 未知类型直接忽略
	}
}

func handleStoryTask(ctx context.Context, task Task) error {
	// TODO: 调用AI生成故事的具体逻辑
	logger.Log.Info("AI story generation placeholder", zap.String("id", task.ID))
	// 模拟处理可能失败的情况
	// return fmt.Errorf("simulated error")
	return nil
}

func handlePoemTask(ctx context.Context, task Task) error {
	logger.Log.Info("AI poem generation placeholder", zap.String("id", task.ID))
	return nil
}
