package main

import (
	"fmt"
)

// inChannel 输入到信道
func inChannel(nums []int) chan int {
	inChan := make(chan int)
	go func() {
		for _, num := range nums {
			inChan <- num
		}
		close(inChan)
	}()
	return inChan
}

// outChannel 输出到信道
func outChannel(inChan <-chan int) chan int {
	finalChan := make(chan int)
	go func() {
		for num := range inChan {
			finalChan <- num * num
		}
		close(finalChan)
	}()
	return finalChan
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	inChan := inChannel(nums)
	finalChan := outChannel(inChan)
	for v := range finalChan {
		fmt.Println(v)
	}
}
