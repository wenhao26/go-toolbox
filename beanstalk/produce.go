package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

func main() {
	conn, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		log.Fatal(err)
	}

	id, _ := conn.Put([]byte("11111111"), 1, 0, 120*time.Second)
	fmt.Println("job", id)
}
