package main

import (
	"fmt"
	"time"
)

type Job struct {
	MsgID string `json:"msg_id"`
}

func main() {
	jobCh := make(chan Job)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case job := <-jobCh:
				fmt.Println(job)
			case <-quit:
				return
			}
		}
	}()

	for i := 0; i < 100; i++ {
		jobCh <- Job{MsgID: time.Now().String()}
	}

	close(jobCh)
	quit <- true
}
