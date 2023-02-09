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
			// 为了防止某个consume长时间占用，设置了timeout
			id, body, err := conn.Reserve(5 * time.Second)
			if err != nil {
				log.Println(err)
				continue
			}
			if id > 0 {
				log.Println("SUB#JOB-ID=", id, ";BODY=", string(body))
				_ = conn.Delete(id)
			}
		}
	}()
	defer wg.Done()
	wg.Wait()
}
