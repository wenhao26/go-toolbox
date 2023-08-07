package main

import (
	"fmt"
)

func process(ch chan int) {
	//time.Sleep(1e9)
	ch <- 1
}

func main() {
	channels := make([]chan int, 100)

	for i := 0; i < 100; i++ {
		channels[i] = make(chan int)
		go process(channels[i])
	}

	for i, ch := range channels {
		<-ch

		fmt.Println("Routine ", i, " quit!")
	}

}
