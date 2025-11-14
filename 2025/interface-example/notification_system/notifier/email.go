package notifier

import (
	"fmt"
)

// EmailNotifier 邮件通知结构体
type EmailNotifier struct {
	SMTPServer string // 邮件发送所需的服务器地址等配置
}

// Send 实现了 Notifier 接口，执行邮件发送的具体逻辑
func (e *EmailNotifier) Send(message string, recipient string) error {
	// 伪代码：实际项目中这里会调用邮件 SDK 或 API，进行网络通信。
	fmt.Printf("[邮件通知] 通过服务器 %s 向 %s 发送邮件：'%s'\n", e.SMTPServer, recipient, message)
	return nil
}
