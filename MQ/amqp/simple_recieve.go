package main

import (
	"toolbox/MQ/rabbitmq"
)

func main() {
	mq := rabbitmq.NewSimple("test.queue_02")
	mq.ConsumerSimple()
}
