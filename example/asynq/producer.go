package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"toolbox/example/asynq/tasks"
)

func main() {
	client := asynq.NewClient(tasks.RedisOpt())
	defer client.Close()

	go func() {
		for {
			task := tasks.CreateTask(tasks.DelayTask, tasks.Task{
				Uid:  999,
				Msg:  tasks.GetMsg(tasks.DelayTask),
				Date: time.Now().String(),
			})
			info, err := client.Enqueue(task, asynq.ProcessIn(10e9))
			if err != nil {
				log.Panic(err.Error())
			}
			fmt.Println("延迟任务=", info.ID)

			time.Sleep(600 * time.Millisecond)
		}
	}()

	for {
		task := tasks.CreateTask(tasks.CommonTask, tasks.Task{
			Uid:  888,
			Msg:  tasks.GetMsg(tasks.CommonTask),
			Date: time.Now().String(),
		})
		info, err := client.Enqueue(task)
		if err != nil {
			log.Panic(err.Error())
		}
		fmt.Println("普通任务=", info.ID)

		time.Sleep(200 * time.Millisecond)
	}
}
