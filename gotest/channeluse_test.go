package gotest

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func receiveTask(ch chan int) {
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("ch close！")
				return
			}
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("task %d is done\n", v)
		}
	}
}

func sendTask() {
	ch := make(chan int, 10)

	go receiveTask(ch)
	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch)
}

func TestCheckTask(t *testing.T) {
	t.Log(runtime.NumGoroutine())
	sendTask()
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(runtime.NumGoroutine())
}
