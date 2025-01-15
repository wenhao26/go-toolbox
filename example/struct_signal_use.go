package main

import (
	"fmt"
	"time"
)

func mockHandle() {
	fmt.Println("wait for 3 sec...")
	time.Sleep(3 * time.Second)
}

func main() {
	done := make(chan struct{})

	go func() {
		mockHandle()
		done <- struct{}{}
		fmt.Println("Done!")
	}()

	<-done
}
