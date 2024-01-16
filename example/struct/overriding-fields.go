package main

import (
	"fmt"
	"time"
)

type Event struct {
	Name      string
	ClickId   string
	PixelId   string
	Timestamp int64
}

type ReportData struct {
	Event
	Properties map[string]interface{}
}

func main() {
	data := ReportData{
		Event: Event{
			Name:      "EVENT_PURCHASE",
			ClickId:   "QdKi0cmhKtQeDpJg/R+xHw==",
			PixelId:   "123456789",
			Timestamp: time.Now().Unix(),
		},
		Properties: map[string]interface{}{
			"event_timestamp": string(time.Now().Unix()),
			"content_type":    "product",
			"currency":        "USD",
			"value":           0.99,
		},
	}
	fmt.Println(data)
}
