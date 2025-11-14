package notifier

import (
	"fmt"
)

// SMSNotifier 实现了 Notifier 接口，执行短信发送的具体逻辑。
type SMSNotifier struct {
	APIToken string // 短信网关所需的 API 令牌等配置
}

// Send 实现了 Notifier 接口，执行邮件发送的具体逻辑
func (s *SMSNotifier) Send(message string, recipient string) error {
	// 伪代码：实际项目中这里会调用短信服务提供商的 API。
	fmt.Printf("[短信通知] 通过 API Token %s 向 %s 发送短信：'%s'\n", s.APIToken, recipient, message)
	return nil
}
