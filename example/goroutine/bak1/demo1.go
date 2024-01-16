package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	ch := make(chan int, 10)

	wg.Add(1)
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Println("向管道ch写入:", i)
		}
		close(ch)
	}()
	go func() {
		defer func() {
			wg.Done()
		}()
		for number := range ch {
			fmt.Println("读取管道ch值：", number)
			fmt.Println("当前通道元素个数", len(ch))
		}
	}()
	wg.Wait()
}
