// Package model 定义生产者/消费者共享的数据结构
package model

import (
	"time"
)

// MessageType 消息的业务类型，不同类型消息在消费者端会有不同的处理耗时
// 用来模拟真实业务场景中“轻量请求”与“重量级计算”混合的负载特征
type MessageType int

const (
	// TypeLight 轻量消息，例如：心跳、埋点上报，处理速度快
	TypeLight MessageType = iota
	// TypeMedium 中等消息，例如：订单状态变更，涉及少量计算或IO
	TypeMedium
	// TypeHeavy 重量级消息，例如：图片/视频转码、大数据聚合，处理耗时较长
	TypeHeavy
)

// Message 是在生产者与消费者之间流转的最小数据单元
type Message struct {
	ID         uint64      // 全局唯一自增ID
	ProducerID int         // 产生该消息的生产者编号，用于追踪数据来源
	Type       MessageType // 消息类型，决定消费端的模拟处理耗时
	Payload    string      // 消息负载，真实场景中可替换为 []byte 或具体业务结构体
	CreatedAt  time.Time   // 消息创建时间，用于计算端到端延迟
}

// String 实现 fmt.Stringer 便于日志打印
func (t MessageType) String() string {
	switch t {
	case TypeLight:
		return "LIGHT"
	case TypeMedium:
		return "MEDIUM"
	case TypeHeavy:
		return "HEAVY"
	default:
		return "UNKNOWN"
	}
}
