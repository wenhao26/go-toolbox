package main

import (
	"fmt"
	"toolbox/2025/interface-example/payment_system/gateway"
	"toolbox/2025/interface-example/payment_system/service"
)

func main() {
	fmt.Println("--- 支付系统启动 ---")

	// 支付宝服务
	alipay := &gateway.AlipayGateway{AppID: "M123456789", Secret: "KAbPYuxc8xqVtwBMjEEH"}
	alipayService := service.PaymentService{Gateway: alipay}
	_, err := alipayService.ProcessPayment("T123456789", "P_$1.99", 1.99)
	if err != nil {
		fmt.Println(err)
	}

}
