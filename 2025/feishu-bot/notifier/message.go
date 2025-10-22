package notifier

import (
	"encoding/json"
)

// Message 消息结构体
type Message struct {
	Payload map[string]interface{}
}

// NewMessage 创建消息实例
func NewMessage() *Message {
	return &Message{Payload: make(map[string]interface{})}
}

// buildPayload 构建载荷数据
func (m *Message) buildPayload(msgType string, data interface{}) *Message {
	m.Payload["msg_type"] = msgType

	dataBytes, _ := json.Marshal(data)

	var dataMap map[string]interface{}
	_ = json.Unmarshal(dataBytes, &dataMap)

	for key, value := range dataMap {
		m.Payload[key] = value
	}

	return m
}

// TextContent 文本消息内容
type TextContent struct {
	Text string `json:"text"`
}

// TextData 文本消息体
type TextData struct {
	Content TextContent `json:"content"`
}

// Text 文本消息
func Text(content string) *Message {
	data := TextData{
		Content: TextContent{
			Text: content,
		},
	}

	return NewMessage().buildPayload("text", data)
}
