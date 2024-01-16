package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{})

	go func() {
		fmt.Println("子协程执行完毕1")
	}()
	go func() {
		fmt.Println("子协程执行完毕2")
	}()
	go func() {
		fmt.Println("子协程执行完毕3")
		done <- struct{}{}
	}()
	<-done

	fmt.Println("主协程执行完毕")

}
