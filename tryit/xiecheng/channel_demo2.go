package main

import (
	"fmt"
	"time"
)

func main() {
	/*ch := make(chan int)
	defer close(ch)

	go func() {
		ch <- 55
	}()
	out := <-ch

	fmt.Println(out)*/
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}

	fmt.Println("End")
}
