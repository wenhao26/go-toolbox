package main

import (
	"fmt"
	"time"
)

func write(ch chan<- string) {
	defer close(ch)
	for i := 0; i < 10; i++ {
		time.Sleep(1e9)
		ch <- time.Now().String()
	}
}

func main() {
	ch := make(chan string)

	go write(ch)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println(v)
		}
	}
}
