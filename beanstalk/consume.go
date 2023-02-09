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

	for {
		id, body, err := conn.Reserve(500 * time.Second)
		if err != nil {
			log.Fatal(err)
		}
		if id > 0 {
			fmt.Println(id, string(body))

			conn.Delete(id)
		}

	}

}
