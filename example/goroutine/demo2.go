package main

import (
	"context"
	"fmt"
)

func main() {
	/*quit := make(chan bool)

	go func() {
		defer fmt.Println("子协程退出")

		for {
			select {
			case <-quit:
				return
			default:
				fmt.Println("Running...")
			}
		}
	}()

	time.Sleep(10 * time.Second)
	quit <- true*/

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		defer fmt.Println("子协程退出")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("Running...")
			}
		}
	}(ctx)
	cancel()

	fmt.Println("主协程退出")
}
