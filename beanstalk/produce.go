package main

import (
	"log"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

func main() {
	conn, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			msg := time.Now().String()
			// 参数分别为任务，优先级，延时时间，处理任务时间
			id, _ := conn.Put([]byte(msg), 1, 0, 10*time.Second)
			log.Println("PUB#JOB-ID=", id)
			time.Sleep(time.Millisecond * 50)
		}
	}()
	defer wg.Done()
	wg.Wait()
}
