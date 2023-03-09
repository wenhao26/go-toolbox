package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	CommonTask = "common:task"
	DelayTask  = "delay:task"
)

// 任务结构体
type Task struct {
	Uid  int    `json:"uid"`
	Msg  string `json:"msg"`
	Date string `json:"date"`
}

// Redis连接配置
func RedisOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
}

// 创建任务
func CreateTask(typename string, t Task) *asynq.Task {
	data, _ := json.Marshal(t)
	return asynq.NewTask(typename, data)
}

func GetMsg(typename string) string {
	switch typename {
	case "common:task":
		return "这是一个普通任务"
	case "delay:task":
		return "这是一个延迟任务"
	default:
		return ""
	}
}
