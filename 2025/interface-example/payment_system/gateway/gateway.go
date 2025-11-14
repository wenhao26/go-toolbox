package gateway

// 定义整个支付系统的行为标准（接口）

// PaymentGateway 定义所有支付网关必须提供的能力
// 抽象层，任何业务逻辑只应该依赖于它
type PaymentGateway interface {
	Charge(orderID string, productID string, amount float64) (bool, error)
}
