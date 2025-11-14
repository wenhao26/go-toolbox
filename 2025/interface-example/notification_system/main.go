package main

import (
	"fmt"
	"toolbox/2025/interface-example/notification_system/notifier"
	"toolbox/2025/interface-example/notification_system/user"
)

func main() {
	fmt.Println("--- 通知系统多态性演示 ---")

	// ----------------------------------------------------
	// 场景 A: 使用 EmailNotifier 实现
	// ----------------------------------------------------

	// 1. 实例化具体的 Email 实现
	emailGateway := &notifier.EmailNotifier{SMTPServer: "smtp.company.com"}

	// 2. 依赖注入：将 Email 实现注入到 UserService 中
	// 赋值是合法的，因为 *EmailNotifier 实现了 notifier.Notifier 接口
	emailUserService := user.UserService{Notifier: emailGateway}

	// 3. 运行业务逻辑
	err := emailUserService.RegisterUser("张三", "zhangsan@mail.com", "138xxxxxxxx")
	if err != nil {
		return
	}

	// ----------------------------------------------------
	// 场景 B: 切换到 SMSNotifier 实现 (核心业务逻辑不变)
	// ----------------------------------------------------

	// 1. 实例化具体的 SMS 实现
	smsGateway := &notifier.SMSNotifier{APIToken: "TOKEN-XYZ-456"}

	// 2. 依赖注入：将 SMS 实现注入到另一个 UserService 实例中
	// 实现了多态切换，user.UserService 的代码无需修改
	smsUserService := user.UserService{Notifier: smsGateway}

	// 3. 运行相同的业务逻辑
	err = smsUserService.RegisterUser("李四", "lisi@mail.com", "139xxxxxxxx")
	if err != nil {
		return
	}
}
