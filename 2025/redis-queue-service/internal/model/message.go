package model

import (
	"encoding/json"
	"time"
)

// ExampleMessage 定义队列消息结构
type ExampleMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ToBytes 将ExampleMessage序列化为JSON字节数组
func (m *ExampleMessage) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

// FromBytes 将字节数组反序列化为ExampleMessage。
func FromBytes(data []byte) (*ExampleMessage, error) {
	var m ExampleMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
