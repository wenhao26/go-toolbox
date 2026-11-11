package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d started\n", id)
	time.Sleep(2 * time.Second)
	fmt.Printf("Worker %d finished\n", id)
}

func main() {
	maxConcurrency := 10
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		sem <- struct{}{} // 主线程尝试获取令牌，卡主循环，防止创建多余协程

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }() // 任务完成后，释放令牌
			worker(id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines finished")
}
