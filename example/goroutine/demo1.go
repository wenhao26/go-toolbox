package main

import (
	"fmt"
	"time"
)

func main() {
	stopCh := make(chan struct{})

	go func() {
		fmt.Println("1111")
		//time.Sleep(5e9)
		close(stopCh)
	}()

	go func() {
		x, ok := <-stopCh
		fmt.Println("2222", x, ok)
	}()

	time.Sleep(5e9)
	fmt.Println("main")
}
