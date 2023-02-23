package main

import (
	"time"

	"toolbox/MQ/rabbitmq"
)

func main() {
	mq := rabbitmq.NewSub("sub-exchange")
	for {
		msg := time.Now().String()
		mq.PublishPub(msg)
		time.Sleep(50 * time.Millisecond)
	}
}
