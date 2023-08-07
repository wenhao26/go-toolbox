package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	{
		go func() {
			time.Sleep(5e9)
			log.Println("Goroutine 1 finished!")
			wg.Done()
		}()
		go func() {
			time.Sleep(2e9)
			log.Println("Goroutine 2 finished!")
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("All Goroutine finished!")
}
