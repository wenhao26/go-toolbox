package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// 定时发送心跳
func sendHeartbeat(ch *amqp.Channel, queueName string) {
	ticker := time.NewTicker(1 * time.Minute) // 每1分钟发送一次心跳
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 发送心跳消息
			err := ch.Publish(
				"",
				queueName,
				false, // 是否持久化
				false, //
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte("heartbeat"),
				},
			)
			if err != nil {
				log.Printf("发送心跳失败: %s", err)
			}
		}
	}
}
