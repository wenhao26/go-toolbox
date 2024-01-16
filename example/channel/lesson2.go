package main

import (
	"fmt"
	"time"
)

func lesson2Task1() <-chan int {
	c := make(chan int, 1)
	go func() {
		time.Sleep(time.Millisecond * 300)
		c <- 1
	}()
	return c
}

func lesson2Task2() <-chan int {
	c := make(chan int, 1)
	go func() {
		time.Sleep(time.Millisecond * 200)
		c <- 2
	}()
	return c
}

// select 多路选择和 time.After 超时
func main() {
	select {
	case val, ok := <-lesson2Task1():
		if !ok {
			fmt.Println("管道1已关闭")
			break
		}
		fmt.Println("管道1：", val)
	case val, ok := <-lesson2Task2():
		if !ok {
			fmt.Println("管道2已关闭")
			break
		}
		fmt.Println("管道2：", val)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout")
	}
	//time.Sleep(1e9)
}
