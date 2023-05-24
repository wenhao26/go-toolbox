package main

import (
	"fmt"
	"log"
	"sync"

	poolService "gitee.com/tym_hmm/rabbitmq-pool-go"
)

var oncePool2 sync.Once
var consumerPool *poolService.RabbitPool

//var consumerPool *poolService.RabbitPool

func initMQ2() *poolService.RabbitPool {
	oncePool2.Do(func() {
		// 初始化消费者
		consumerPool = poolService.NewConsumePool()
		consumerPool.SetMaxConnection(100)
		err := consumerPool.Connect("localhost", 5672, "admin", "admin")
		if err != nil {
			log.Fatalln(err)
		}
	})
	return consumerPool
}

func main() {
	initMQ2()

	cr := &poolService.ConsumeReceive{
		ExchangeName: "test-exchange1",
		ExchangeType: poolService.EXCHANGE_TYPE_DIRECT,
		Route:        "test-route1",
		QueueName:    "test-queue1",
		IsTry:        true,
		IsAutoAck:    false,
		MaxReTry:     5,
		EventFail: func(code int, e error, data []byte) {
			fmt.Printf("error:%s", e)
		},
		EventSuccess: func(data []byte, header map[string]interface{}, retryClient poolService.RetryClientInterface) bool {
			_ = retryClient.Ack()
			fmt.Printf("data:%s\n", string(data))
			return true
		},
	}
	consumerPool.RegisterConsumeReceive(cr)
	err := consumerPool.RunConsume()
	if err != nil {
		fmt.Println(err)
	}
}
