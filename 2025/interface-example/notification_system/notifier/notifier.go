package notifier

// Notifier 接口：定义了所有通知渠道必须提供的行为标准
type Notifier interface {
	Send(message string, recipient string) error
}
