package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func Task(n, runID int) {
	for {
		if n > 1000 {
			fmt.Printf("--RunID=%d Done...\n", runID)
			break
		}
		n++
		fmt.Printf(" --RunID=%d;This`s number=%d\n", runID, n+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup

	p, _ := ants.NewPool(20, ants.WithPreAlloc(false))
	defer p.Release()

	go func(p *ants.Pool) {
		for {
			fmt.Printf(" ##Running Pool number:%d\n", p.Running())
			if p.Running() == 0 {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}(p)

	runNumber := 10
	for i := 0; i < runNumber; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			n := rand.Intn(100)
			Task(n, i)
			wg.Done()
		})
		fmt.Printf("ID=%d;Running Pool number:%d\n", i, p.Running())
	}
	wg.Wait()
}
