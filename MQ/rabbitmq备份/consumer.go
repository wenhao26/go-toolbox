package rabbitmq备份

import (
	"github.com/streadway/amqp"
)

// Consumer 定义消费者选项
type Consumer struct {
	Name      string
	AutoAck   bool // 自动确认
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func defaultConsumer() *Consumer {
	return &Consumer{"", true, false, false, false, nil}
}
