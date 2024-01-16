package main

import (
	"fmt"
	"sync"
)

// 使用两个goroutine交替打印序列，一个gorountine打印数字，另一个goroutine打印字母,最终效果如下：
// 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
func main() {
	number := make(chan bool)
	letter := make(chan bool)

	var wg sync.WaitGroup
	go func() {
		i := 1
		for {
			select {
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- true
			}
		}
	}()

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		i := 'A'
		for {
			select {
			case <-letter:
				if i >= 'Z' {
					wg.Done()
					return
				}
				fmt.Print(string(i))
				i++
				fmt.Print(string(i))
				i++
				number <- true
			}
		}
	}(&wg)

	number <- true
	wg.Wait()
}
