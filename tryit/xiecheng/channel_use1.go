package main

import (
	"fmt"
	"time"
)

func main() {
	Test()
}

func Run(f func(s string, ch chan string), s string, timeout int, cOut chan string) {
	ch_run := make(chan string)
	// go run(s, ch_run)
	go f(s, ch_run)

	select {
	case re := <-ch_run:
		cOut <- re
	case <-time.After(time.Duration(timeout) * time.Second):
		re := fmt.Sprintf("task timeout ")
		cOut <- re
	}
}

// func run(s string, ch chan string) {
//  time.Sleep(time.Duration(3) * time.Second)

//  ch <- fmt.Sprintf("task input %s,sleep %d second", s, 3)
//  return
// }

func aa1(s string, ch chan string) {
	time.Sleep(time.Duration(3) * time.Second)
	ch <- fmt.Sprintf("task1 input %s,sleep %d second", s, 3)

}

func aa2(s string, ch chan string) {
	time.Sleep(time.Duration(5) * time.Second)
	ch <- fmt.Sprintf("task2 input %s,sleep %d second", s, 5)

}

func aa3(s string, ch chan string) {
	time.Sleep(time.Duration(10) * time.Second)
	ch <- fmt.Sprintf("task3 input %s,sleep %d second", s, 10)

}

func Test() {
	a := synchron(20, "aaa", aa1, aa2, aa3)
	fmt.Printf("result: %v \n", a)
}

// timeout: 超时时间
// input: 统一入参
// args: 方法
func synchron(timeout int, input string, args ...func(s string, ch chan string)) []string {
	// input := []string{"aaa", "bbb", "ccc"}
	// timeout := 8
	// 创建N个任务管道，用来接收各个并发任务的完成结果
	chs := make([]chan string, len(args))

	defer func() {
		for _, c := range chs {
			if c != nil {
				close(c)
			}
		}
	}()

	sTime := time.Now()
	fmt.Println("start")

	for i, f := range args {
		chs[i] = make(chan string)
		go Run(f, input, timeout, chs[i])
	}

	resList := []string{}
	// 获取结果
	for _, ch := range chs {
		v := <-ch
		fmt.Println(v)
		resList = append(resList, v)
	}

	eTime := time.Now()
	fmt.Printf("finished,Process time %s. Number of task is %d \n", eTime.Sub(sTime), len(args))
	// 将多个异步任务同时返回
	return resList
}
