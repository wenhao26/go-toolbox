// 实现一个限流器，限制 Goroutine 的并发数量。
package main

import (
	"fmt"
	"sync"
	"time"
)

func work(num int, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	sem <- struct{}{}        // 获取信号
	defer func() { <-sem }() // 释放信号

	fmt.Printf("Worker %d started\n", num)
	time.Sleep(2 * time.Second)
	fmt.Printf("Worker %d finished\n", num)
}

func main() {
	maxConcurrency := 3
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		work(i, sem, &wg)
	}

	wg.Wait()

	fmt.Println("END")
}
