package main

import (
	"fmt"
	"time"
)

func main() {
	exit := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("output:", i)
		}
		exit <- true
	}()

	<-exit
	fmt.Println("end")
}
