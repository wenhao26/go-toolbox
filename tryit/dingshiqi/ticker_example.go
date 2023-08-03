package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"toolbox/utils"
)

func tickerDemo1() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running...")
	}

}

// 模拟乘客
func GetNewPassenger() string {
	time.Sleep(2e9)

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(2)
	if n == 1 {
		return utils.GenUUID()
	}
	return ""
}

// 模拟发车
func drive(passengers []string) {
	log.Printf(">>>>>Time's up, let's go!!! Total passengers=%d \r\n", len(passengers))
}

// 模拟公交车运行
// 1、每隔一分钟一班次
// 2、如果满员，不用等到开车时间，直接开走
func tickerDemo2() {
	ticker := time.NewTicker(1 * time.Minute)
	log.Println("模拟每隔一分钟，发出一趟大巴车")

	// 总班次
	maxShift := 3
	totalShift := 0

	// 最大装载乘客人数
	maxPassenger := 30
	passengers := make([]string, 0, maxPassenger)

	for {
		if totalShift > maxShift {
			fmt.Println("come or go off wor!!!")
			ticker.Stop()
			break
		}

		fmt.Println("Waiting...")

		passenger := GetNewPassenger()
		if passenger != "" {
			log.Printf("Passenger boarding:%s", passenger)
			passengers = append(passengers, passenger)
		} else {
			// 模拟继续“候车”
			time.Sleep(1e9)
		}

		select {
		case <-ticker.C: // 到时间了，发车
			totalShift++
			drive(passengers)
			passengers = []string{}
		default:
			if len(passengers) >= maxPassenger { // 时间没到，但是乘客满员了，发车
				totalShift++
				drive(passengers)
				passengers = []string{}
			}
		}
	}

}

func main() {
	//tickerDemo1()
	tickerDemo2()
}
