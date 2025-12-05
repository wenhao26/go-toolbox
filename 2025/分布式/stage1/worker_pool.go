// 并发工作池，工作池用于限制并发执行的任务数量，有效地管理系统资源
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	taskNum   = 9
	workerNum = 3
)

func worker(id int, tasks <-chan int, results chan<- int) {
	for task := range tasks {
		fmt.Printf("Worker %d 开始处理任务 %d...\n", id, task)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d 完成任务 %d。\n", id, task)
		results <- task * 2
	}
}

func main() {
	tasks := make(chan int, taskNum)
	results := make(chan int, taskNum)

	var wg sync.WaitGroup
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, tasks, results)
		}(i)
	}

	for i := 1; i <= taskNum; i++ {
		tasks <- i
	}
	close(tasks)

	wg.Wait()

	fmt.Println("\n所有任务和结果：")
	for num := 1; num <= taskNum; num++ {
		result := <-results
		fmt.Printf("任务 %d 的结果是 %d\n", num, result)
	}
}
