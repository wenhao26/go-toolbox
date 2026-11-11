package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func handleJob(id int, name string, resultCh chan<- string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Task %d panic: %v\n", id, r)
		}
	}()

	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	str := fmt.Sprintf("📑 Worker %d completed task: %s", id, name)
	fmt.Println(str)

	resultCh <- str
}

func main() {
	jobs := []string{"T0", "T1", "T2", "T3", "T4", "T5", "T6", "T7", "T8", "T9"}
	resultCh := make(chan string, len(jobs))

	pool, err := ants.NewPool(5)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	var wg sync.WaitGroup

	for i, job := range jobs {
		id := i + 1
		name := job

		wg.Add(1)
		err := pool.Submit(func() {
			defer wg.Done()
			handleJob(id, name, resultCh)
		})
		if err != nil {
			fmt.Println(err)
			wg.Done()
		}
	}

	go func() {
		wg.Wait()
		close(resultCh)
		fmt.Println("Closed!")
	}()

	for result := range resultCh {
		fmt.Println("✅ Result:", result)
	}
}
