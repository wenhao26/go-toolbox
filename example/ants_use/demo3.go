package main

import (
	"fmt"
	"sync"
	"time"
)

func read(ch chan bool, wg *sync.WaitGroup, n int) {
	defer wg.Done()

	ch <- true
	fmt.Printf("number：%d, time：%s \n\r", n, time.Now().String())
	<-ch
}

func main() {
	var wg sync.WaitGroup
	count := 10

	ch := make(chan bool, 2)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go read(ch, &wg, i)
	}
	wg.Wait()
}
