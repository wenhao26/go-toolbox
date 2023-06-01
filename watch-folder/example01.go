package main

import (
	"fmt"
	"time"
)

func main() {
	msgCh := make(chan string, 100)

	send(msgCh)
	receive(msgCh)

	select {}
}

func send(msgCh chan<- string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			date := time.Now().String()
			msgCh <- date
			fmt.Println("SEND-", date)
			time.Sleep(1e9)
		}
	}()
}

func receive(msgCh <-chan string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			select {
			case data, ok := <-msgCh:
				if !ok {
					return
				}
				fmt.Println("RECEIVE-", data)
			}
		}
	}()
}
