package main

import (
	"fmt"
	"log"

	"toolbox/2025/order-example/order"
)

func main() {
	userID := "uid_1205"
	productID := "book_sn95270001245"
	quantity := 2
	unitPrice := 18.90
	contactNumber := "xxx-123456789"
	deliveryAddress := "蓝色星球-兔子省-茶叶市-什么都不想说街道-向东西南北001-99号"

	o := order.NewOrder(userID, productID)
	//orderID := o.GenerateOrderID()
	//fmt.Println(orderID)

	//Create(userID, productID string, quantity int, unitPrice float32, contactNumber, deliveryAddress string)
	createItem, err := o.Create(userID, productID, quantity, unitPrice, contactNumber, deliveryAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(createItem)

	getItem := o.Get(createItem.OrderID)
	fmt.Println(getItem)

	updateItem := o.Update(getItem.OrderID, order.Unpaid)
	fmt.Println(updateItem)
}
