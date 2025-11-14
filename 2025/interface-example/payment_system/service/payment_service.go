package service

import (
	"fmt"
	"toolbox/2025/interface-example/payment_system/gateway"
)

// PaymentService z核心业务服务
type PaymentService struct {
	Gateway gateway.PaymentGateway // 依赖于 PaymentGateway 接口
}

// ProcessPayment 核心逻辑：稳定且独立于支付渠道
func (s *PaymentService) ProcessPayment(orderID string, productID string, amount float64) (bool, error) {
	fmt.Printf("\n--- 核心业务逻辑启动: 订单 %s ---\n", orderID)

	// 业务校验...

	// 调用抽象方法：利用多态性
	_, err := s.Gateway.Charge(orderID, productID, amount)

	// 结果处理...
	if err != nil {
		fmt.Printf("核心逻辑: 支付渠道调用失败: %v\n", err)
		return false, err
	}
	// ...
	fmt.Printf("核心逻辑: 订单 %s 状态更新为“已支付”。\n", orderID)
	return true, nil
}
