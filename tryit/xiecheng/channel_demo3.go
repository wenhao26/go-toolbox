package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)
	quit := make(chan struct{})

	rand.Seed(time.Now().UnixNano())

	go func() {
		jobNum := rand.Intn(10)
		fmt.Println("jobs", jobNum)

		wg.Add(jobNum)
		for i := 0; i < jobNum; i++ {
			go func() {
				defer wg.Done()
				time.Sleep(2e9)
				ch <- time.Now().String()
			}()
		}
		wg.Wait()
		close(quit)
	}()

	for {
		select {
		case r := <-ch:
			fmt.Println("from ch:", r)
		case <-quit:
			fmt.Println("exit！")
			return
		}
	}
}
