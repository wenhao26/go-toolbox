package main

import (
	"log"
	"runtime"
	"time"
)

func DoSomething() {
	for {
		// todo
		log.Println(time.Now().String())
		runtime.Gosched()
	}
}

func main() {
	go DoSomething()
	go DoSomething()
	select {}
}
