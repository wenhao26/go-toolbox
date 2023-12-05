package main

import (
	"fmt"
	"time"
)

func main() {
	runCount := 0
	ticker := time.NewTicker(3e9)

	for {
		if runCount > 10 {
			ticker.Stop()
			break
		}

		select {
		case <-ticker.C:
			fmt.Println("Running...")
			runCount++
		default:
			//fmt.Println("...")
		}
	}

}
