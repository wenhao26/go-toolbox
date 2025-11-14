package user

import (
	"fmt"
	"toolbox/2025/interface-example/notification_system/notifier"
)

// UserService 用户服务：包含核心业务流程
type UserService struct {
	Notifier notifier.Notifier
}

// RegisterUser 注册用户业务逻辑
func (u *UserService) RegisterUser(username string, email string, phone string) error {
	fmt.Printf("\n--- 业务流程：用户 %s 正在注册 ---\n", username)

	// 1. 业务逻辑：验证数据、存储到数据库...

	// 2. 发送通知 (核心多态体现)：
	// 无论 s.Notifier 实际是 Email 还是 SMS，我们都调用 Send 方法。
	// Go 运行时会自动执行正确实现的代码。

	// 假设通知发送到用户的邮箱/手机号
	notificationMsg := fmt.Sprintf("欢迎您，%s！注册成功。", username)

	// 这里通过接口调用 Send，实现了多态
	err := u.Notifier.Send(notificationMsg, email)

	if err != nil {
		fmt.Printf("业务日志：通知发送失败，但用户注册可能已成功：%v\n", err)
	}

	return nil
}
