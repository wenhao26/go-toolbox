package main

import (
	"fmt"
	"sync"
)

func MyWorker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		results <- job * 2
	}
	wg.Done()
}

func main() {
	numJobs := 10
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go MyWorker(i, jobs, results, &wg)
	}

	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println("Result:", result)
	}
}
