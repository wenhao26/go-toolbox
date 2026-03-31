package service

import (
	"context"
	"encoding/json"
	"fmt"
	"toolbox/2026/gocron/internal/model"

	"github.com/redis/go-redis/v9"
)

// -- 任务模型与存储层 --

var (
	RDB *redis.Client
	ctx = context.Background()
)

const (
	KeyTaskMap = "gocron:tasks"   // Hash存储所有任务元数据
	KeyLogList = "gocron:logs:%s" // List存储每个任务的日志
)

// InitRedis 初始化 Redis
func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{Addr: addr})
}

// SaveTask 保存任务
func SaveTask(task model.Task) error {
	data, _ := json.Marshal(task)
	return RDB.HSet(ctx, KeyTaskMap, task.ID, data).Err()
}

// DeleteTask 删除任务
func DeleteTask(id string) {
	RDB.HDel(ctx, KeyTaskMap, id)
	RDB.Del(ctx, fmt.Sprintf(KeyLogList, id))
}

// LoadAllTasks 加载所有任务
func LoadAllTasks() ([]model.Task, error) {
	results, err := RDB.HGetAll(ctx, KeyTaskMap).Result()
	if err != nil {
		return nil, err
	}

	var tasks []model.Task

	for _, v := range results {
		var task model.Task
		_ = json.Unmarshal([]byte(v), &task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// SaveLog 保存日志并保留最新50条
func SaveLog(taskID string, log model.Log) {
	data, _ := json.Marshal(log)
	key := fmt.Sprintf(KeyLogList, taskID)
	RDB.LPush(ctx, key, data)
	RDB.LTrim(ctx, key, 0, 49)
}

// GetLogs 获取历史日志
func GetLogs(taskID string) []model.Log {
	key := fmt.Sprintf(KeyLogList, taskID)
	results, _ := RDB.LRange(ctx, key, 0, 19).Result()
	var logs []model.Log
	for _, v := range results {
		var log model.Log
		_ = json.Unmarshal([]byte(v), &log)
		logs = append(logs, log)
	}
	return logs
}
