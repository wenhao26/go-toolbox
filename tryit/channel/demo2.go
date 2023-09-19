package main

import (
	"fmt"
)

// 购买
func purchase(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			out <- fmt.Sprintf("purchase task %d", i)
		}
	}()
	return out
}

// 构建
func build(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			out <- fmt.Sprintf("build...%s", v)
		}
	}()
	return out
}

// 完成
func complete(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			out <- fmt.Sprintf("complete...%s", v)
		}
	}()
	return out
}

func main() {
	p := purchase(10)
	b := build(p)
	c := complete(b)

	for v := range c {
		fmt.Println(v)
	}

}
