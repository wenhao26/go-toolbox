package main

import (
	"fmt"
	"log"
	"sync"

	poolService "gitee.com/tym_hmm/rabbitmq-pool-go"
)

var oncePool sync.Once
var producerPool *poolService.RabbitPool
//var consumerPool *poolService.RabbitPool

func initMQ() *poolService.RabbitPool {
	oncePool.Do(func() {
		// 初始化生产者
		producerPool = poolService.NewProductPool()
		// 初始化消费者
		//consumerPool = poolService.NewConsumePool()

		err := producerPool.Connect("localhost", 5672, "admin", "admin")
		if err != nil {
			log.Fatalln(err)
		}
	})
	return producerPool
}

func main() {
	initMQ()

	var wg sync.WaitGroup
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			data := poolService.GetRabbitMqDataFormat(
				"test-exchange1",
				poolService.EXCHANGE_TYPE_DIRECT,
				"test-queue1",
				"test-route1",
				fmt.Sprintf("data=%d", num),
				)
			err := producerPool.Push(data)
			if err != nil {
				fmt.Println(err)
			}
		}(i)
	}
	wg.Wait()
}
