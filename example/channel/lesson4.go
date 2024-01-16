package main

import (
	"context"
	"fmt"
	"time"
)

func firstCtx(ctx context.Context) {
	go secondCtx(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("first done")
			return
		default:
			fmt.Println("first running")
			time.Sleep(2e9)
		}
	}
}

func secondCtx(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("second done")
			return
		default:
			fmt.Println("second running")
			time.Sleep(2e9)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go firstCtx(ctx)
	time.Sleep(5e9)

	fmt.Println("stop all sub goroutine")
	cancel()
	time.Sleep(5e9)
}
