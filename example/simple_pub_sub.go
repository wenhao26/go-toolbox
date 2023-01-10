package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 生产者
func Publisher(msgCh chan<- string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		for {
			curDate := time.Now().String()
			msgCh <- curDate
			fmt.Println("[sent]=", curDate)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

// 订阅者
func Subscriber(msgCh <-chan string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		for {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				fmt.Println("[receive]=", msg)
			}
		}
	}()
}

func main() {
	msgCh := make(chan string, 100)
	defer close(msgCh)

	Publisher(msgCh)
	Subscriber(msgCh)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
