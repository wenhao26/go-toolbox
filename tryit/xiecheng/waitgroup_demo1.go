package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var mutex sync.Mutex
var sum int

func add()  {
	defer wg.Done()

	mutex.Lock()
	sum++
	mutex.Unlock()
}

func main() {
	n := 100000
	
	wg.Add(n)
	for i := 0; i < n; i++ {
		go add()
	}
	wg.Wait()

	fmt.Println(sum)
}
