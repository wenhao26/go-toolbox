package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

func main() {
	beeep.AppName = "计划提示"

	err := beeep.Notify("Go计划任务提醒", "今天周四记得点餐", "")
	if err != nil {
		panic(err)
	}

	fmt.Println("OK")
}
