package queue

import (
	"errors"
)

// 定义队列服务可能返回的自定义错误
var (
	ErrInvalidMessage   = errors.New("无效的消息格式")
	ErrMessageTooLarge  = errors.New("消息过大")
	ErrProcessingFailed = errors.New("消息处理失败")
	ErrQueueClosed      = errors.New("队列已关闭")
)
