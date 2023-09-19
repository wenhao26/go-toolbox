package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(3e9)
		ch <- "server response"
	}()

	select {
	case v := <-ch:
		fmt.Println(v)
		close(ch)
	case <-time.After(5e9):
		fmt.Println("timeout")
		return
	}
}
