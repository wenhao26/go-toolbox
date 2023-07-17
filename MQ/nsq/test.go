package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Println(i)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Stop(sigChan)
	<-sigChan

	fmt.Println("end...")
}
