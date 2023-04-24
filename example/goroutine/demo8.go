package main

import (
	"fmt"
	"time"
)

func run() {
	fmt.Println(time.Now().String())
	time.Sleep(1e9)
}

func main() {
	/*for {
		go run()
	}*/

	timer := time.NewTimer(time.Second * 3)
	for {
		select {
		case <-timer.C: // 超时
			fmt.Println("timeout")
			timer.Reset(time.Second * 3)
		}
	}

}
