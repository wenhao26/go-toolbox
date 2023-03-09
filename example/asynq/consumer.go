package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"toolbox/example/asynq/tasks"
)

func doTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.Task
	_ = json.Unmarshal(t.Payload(), &p)
	fmt.Println(p)
	return nil
}

func main() {
	server := asynq.NewServer(tasks.RedisOpt(), asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()
	go mux.HandleFunc(tasks.CommonTask, doTask)
	go mux.HandleFunc(tasks.DelayTask, doTask)
	if err := server.Run(mux); err != nil {
		panic(err)
	}
}
