package main

import (
	"fmt"
	"runtime"
	"time"

	"toolbox/example/goroutine/case01/t"
)

type Score struct {
	Num int
}

func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	time.Sleep(1e9)
}

func main() {
	// 注册工作池，传输任务
	poolNum := 200000
	p := t.NewWorkerPool(poolNum)
	p.Run()

	// 模拟100万并发场景
	dataNum := 1000000
	go func() {
		for i := 0; i <= dataNum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine():", runtime.NumGoroutine())
		time.Sleep(2e9)
	}
}
