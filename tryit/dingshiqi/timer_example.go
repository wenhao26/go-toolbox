package main

import (
	"log"
	"time"
)

// 延迟执行一个函数
func delayFunc1() {
	timer := time.NewTimer(5 * time.Second)

	select {
	case <-timer.C:
		log.Println("延迟5秒!")
	}
}

func delayFunc2() {
	log.Println("Start...")

	// 延迟1秒之后
	<-time.After(1 * time.Second)

	log.Println("End...")
}

func delayFunc3() {
	log.Println("Start...")

	// 延迟2秒之后，执行回调函数
	time.AfterFunc(2*time.Second, func() {
		log.Println("End...")
	})

	time.Sleep(3 * time.Second) // 等待协程退出
}

func main() {
	//delayFunc1()
	//delayFunc2()
	delayFunc3()
}
