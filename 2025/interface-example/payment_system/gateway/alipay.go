package gateway

import (
	"errors"
	"fmt"
)

// AlipayGateway 支付宝支付网关
type AlipayGateway struct {
	AppID  string
	Secret string
}

// Charge 实现 PaymentGateway 接口
func (a *AlipayGateway) Charge(orderID string, productID string, amount float64) (bool, error) {
	fmt.Printf(
		"[支付宝] 正在使用 AppID: %s 处理订单: %s, 商品: %s, 金额: %.2f\n",
		a.AppID,
		orderID,
		productID,
		amount,
	)
	// ... 实际的支付宝 API 请求逻辑 ...
	if amount > 10000 {
		return false, errors.New("支付宝接口: 单笔金额超限")
	}
	return true, nil
}
