package main

import (
	"fmt"
	"time"
)

func lesson1Task1() int {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("任务1...")
	return 10
}

func lesson1Task2() {
	fmt.Println("任务2...")
	time.Sleep(time.Millisecond * 200)
	fmt.Println("任务2执行结束")
}

func lesson1RunTasks() chan int {
	ch := make(chan int)
	go func() {
		val := lesson1Task1()
		ch <- val
	}()
	return ch
}

// channel在多个goroutine之间进行通信
func main() {
	ch := lesson1RunTasks()
	lesson1Task2()
	fmt.Println(<-ch)
}
