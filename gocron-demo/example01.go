package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// 创建调度器
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	// 添加任务到调度器
	f := func(a string, b int) {
		t := time.Now().String()
		fmt.Println(t, "=>", a, b)
	}

	j, err := s.NewJob(
		gocron.DurationJob(1*time.Second),
		gocron.NewTask(f, "hello", 1688),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(j.ID())

	s.Start()

	select {
	case <-time.After(time.Hour):
	}

	err = s.Shutdown()
	if err != nil {
		panic(err)
	}
}
