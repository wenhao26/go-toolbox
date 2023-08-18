package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 监控
func monitor(ch chan<- int64) {
	defer close(ch)
	for {
		time.Sleep(100 * time.Millisecond)
		t := time.Now().Unix()
		ch <- t
		fmt.Println("monitor:", t)
	}
}

// 处理
func solver(ch <-chan int64) {
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println("solver:", v)
		}
	}
}

func main() {
	ch := make(chan int64)
	go monitor(ch)
	go solver(ch)

	for i := 0; i < 10; i++ {
		time.Sleep(1e9)
		fmt.Println("main:", i)
	}

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT)
	<-s
}
