package main

import (
	"toolbox/MQ/rabbitmq"
)

func main() {
	mq := rabbitmq.NewSimple("test.lscj")
	mq.ConsumerSimple()
}
