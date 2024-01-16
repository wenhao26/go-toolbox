package main

import (
	"fmt"
	"strconv"
	"time"
)

func doSomething(action string) {
	for {
		fmt.Println("并发ID=", action)
		time.Sleep(1e9)
	}
}

func main() {
	/*for i := 0; i < 10; i++ {
		go doSomething("TEST" + strconv.Itoa(i))
	}*/

	for i := 0; i < 1000000; i++ {
		go doSomething("TEST" + strconv.Itoa(i))
	}

	for {
		fmt.Println("Main")
		time.Sleep(1e9)
	}
}
