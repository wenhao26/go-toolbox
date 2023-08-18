package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func demo4Func1() {
	time.Sleep(1e9)
	fmt.Printf("time：%s \n\r", time.Now().String())
}

func main() {
	defer ants.Release()
	var wg sync.WaitGroup

	runTimes := 10
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(func() {
			demo4Func1()
			defer wg.Done()
		})
	}
	wg.Wait()

	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")
}
