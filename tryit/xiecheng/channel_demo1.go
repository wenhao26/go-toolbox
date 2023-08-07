package main

import (
	"fmt"
)

func Count(ch chan int) {
	ch <- 1

	fmt.Println("Running...")
}

func main() {
	numbers := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		numbers[i] = make(chan int)
		go Count(numbers[i])
	}

	for i, ch := range numbers {
		<-ch
		fmt.Println("Routine ", i, " quit!")
	}
}
