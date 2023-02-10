package main

import (
	"runtime"
	"sync"

	"toolbox/beanstalk/bean"
)

func main() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(runtime.NumCPU())

	b := bean.New()
	defer b.CloseBean()

	wg.Add(1)
	go func() {
		for {
			b.Consume("C1", "channel1")
		}
	}()
	wg.Add(1)
	go func() {
		for {
			b.Consume("C2", "channel1")
		}
	}()
	defer wg.Done()
	wg.Wait()
}
