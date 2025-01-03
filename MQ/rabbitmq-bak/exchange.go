package rabbitmq_bak

import (
	"github.com/streadway/amqp"
)

// Exchange 定义交换机
type Exchange struct {
	Name       string // 名称
	Kind       string // 交换机类型
	routingKey string // 交换机绑定时用的key
	Durable    bool   // 是否持久化
	AutoDelete bool   // 是否自动删除
	Internal   bool   // 是否内置，如果设置为true，则表示是内置的交换器，客户端程序无法直接发送消息到这个交换器，只能通过交换器路由到交换器的方式
	NoWait     bool   // 是否等待通知定义交换机结果
	Args       amqp.Table
	IsDeclare  bool // 是否已定义
}

func defaultExchange(name string) *Exchange {
	return &Exchange{
		Name:    name,
		Kind:    amqp.ExchangeDirect,
		Durable: true, // 持久化
	}
}
