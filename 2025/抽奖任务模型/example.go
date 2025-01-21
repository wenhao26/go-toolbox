// ants 实现协程池

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func main() {
	pool, err := ants.NewPool(3)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		taskID := i // 避免闭包捕获循环变量
		_ = pool.Submit(func() {
			defer wg.Done()
			fmt.Printf("任务 %d 开始执行\n", taskID)
			time.Sleep(1 * time.Second) // 模拟任务执行
			fmt.Printf("任务 %d 执行完成\n", taskID)
		})
	}

	wg.Wait()
	fmt.Println("所有任务已完成")
}
