package main

import (
	"time"

	"toolbox/MQ/rabbitmq"
)

func main() {
	mq := rabbitmq.NewSimple("test.queue_02")
	for {
		msg := time.Now().String()
		mq.PublishSimple(msg)
		time.Sleep(50 * time.Millisecond)
	}

}
