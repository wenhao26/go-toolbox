package main

import (
	"fmt"
	"sync"
	"time"
)

func handler(num int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("goroutine-", num)
	time.Sleep(2e9)
}

func main() {
	var wg sync.WaitGroup

	num := 50
	wg.Add(num)
	for i := 0; i <= num; i++ {
		go handler(i, &wg)
	}

	wg.Wait()
	fmt.Println("End...")
}
