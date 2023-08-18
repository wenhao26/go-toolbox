package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("ants demo")
}

func main() {
	defer ants.Release()

	runTimes := 1000

	var wg sync.WaitGroup
	func1 := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(func1)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")

}
