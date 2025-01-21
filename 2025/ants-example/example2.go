//有一批任务需要处理，每个任务的处理时间不同。
//使用 Goroutine 并发处理任务，使用 Channel 传递任务结果。
//主 Goroutine 汇总所有结果并输出。
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func doTask(id int, name string, resultCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	resultCh <- fmt.Sprintf("Worker %d completed task: %s", id, name)
}

func main() {
	tasks := []string{"task0", "task1", "task2", "task3", "task4", "task5", "task6", "task7"}
	resultChan := make(chan string, len(tasks))

	var wg sync.WaitGroup

	pool, err := ants.NewPool(4)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	for i, task := range tasks {
		wg.Add(1)

		taskId := i + 1
		taskName := task
		_ = pool.Submit(func() {
			doTask(taskId, taskName, resultChan, &wg)
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
		fmt.Println("Closed")
	}()

	for result := range resultChan {
		fmt.Println(result)
	}
}
