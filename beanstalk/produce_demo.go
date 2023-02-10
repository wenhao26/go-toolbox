package main

import (
	"runtime"
	"sync"
	"time"

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
			msg := []byte(time.Now().String())
			b.Produce("P1", "channel1", msg, 1, 0, 5*time.Second)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	defer wg.Done()
	wg.Wait()
}
