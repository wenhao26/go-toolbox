package main

import (
	"toolbox/MQ/rabbitmq"
)

func main() {
	mq := rabbitmq.NewSub("sub-exchange")
	mq.ReceiveSub()
}
