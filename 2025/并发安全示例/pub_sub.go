package main

import (
	"fmt"
	"time"
)

func pub(result chan<- string, done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("[PUB]接收到停止信号，协程即将退出")
			return
		default:
			result <- time.Now().Format("2006-01-02 15:04:05")
			time.Sleep(3 * time.Second)
		}
	}
}

func sub(result <-chan string, done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("[SUB]接收到停止信号，协程即将退出")
			return
		default:
			fmt.Println("-DATE:", <-result)
		}
	}
}

func main() {
	result := make(chan string, 10)
	done := make(chan struct{})

	go pub(result, done)
	go sub(result, done)

	for i := 0; i < 6; i++ {
		currentDate := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("-INFO:", currentDate)
		time.Sleep(1 * time.Second)
	}
	close(done) // 关闭通道，通知协程停止

	fmt.Println("发送完成信号")
	close(result)
}
