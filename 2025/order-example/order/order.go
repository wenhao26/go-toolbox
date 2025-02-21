// 封装订单管理：创建、查询、更新等
package order

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// 订单状态
const (
	Ordered  = "ORDERED"  // 已订购
	Unpaid   = "UNPAID"   // 待支付
	Shipped  = "SHIPPED"  // 已发货
	Received = "RECEIVED" // 已签收
)

// Order 订单接口
type Order interface {
	GenerateOrderID(userID, productID string) string
	Create(userID, productID string, quantity int, unitPrice float32, contactNumber, deliveryAddress string) (order, error)
	Get(orderID string) order
	Update(orderID, status string) order
}

// order 订单结构体
type order struct {
	OrderID         string
	Status          string
	Created         string
	UserID          string
	ProductID       string
	Quantity        int
	UnitPrice       float64
	PaymentAmount   float64
	ContactNumber   string
	DeliveryAddress string
}

// NewOrder 创建订单实例
func NewOrder(userID, productID string) *order {
	return &order{
		UserID:    userID,
		ProductID: productID,
	}
}

// GenerateOrderID 生成订单ID
func (o *order) GenerateOrderID() string {
	currentDate := time.Now().Format("20060102")
	timestamp := time.Now().UnixNano()

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1000000)

	// 组合字符串
	uniqueString := fmt.Sprintf("%s-%s-%d-%d", o.UserID, o.ProductID, timestamp, randomNum)

	// 使用 SHA-256 哈希来确保一致性
	hash := sha256.Sum256([]byte(uniqueString))
	hashString := hex.EncodeToString(hash[:])

	orderID := currentDate + hashString[:14]

	return orderID
}

// Create 创建订单
func (o *order) Create(userID, productID string, quantity int, unitPrice float64, contactNumber, deliveryAddress string) (order, error) {
	orderID := o.GenerateOrderID()
	created := time.Now().String()
	paymentAmount := float64(quantity) * unitPrice

	return order{
		OrderID:         orderID,
		Status:          Ordered,
		Created:         created,
		UserID:          userID,
		ProductID:       productID,
		Quantity:        quantity,
		UnitPrice:       unitPrice,
		PaymentAmount:   paymentAmount,
		ContactNumber:   contactNumber,
		DeliveryAddress: deliveryAddress,
	}, nil
}

// Get 查询订单
func (o *order) Get(orderID string) order {
	return order{
		OrderID: orderID,
	}
}

// Update 更新订单信息
func (o *order) Update(orderID, status string) order {
	return order{
		OrderID: orderID,
		Status:  status,
	}
}
