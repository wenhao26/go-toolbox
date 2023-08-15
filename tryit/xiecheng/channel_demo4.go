package main

import (
	"fmt"
	"time"
)

type Device struct {
	brand   string
	version string
}

func receive(ch chan Device)  {
	val := <- ch
	fmt.Println(val)
}

func main() {
	ch := make(chan Device)
	device := Device{brand: "HUAWEI", version: "HarmonyOS 4.0"}

	go receive(ch)
	// 往通道发送数据
	ch <- device

	time.Sleep(2e9)
	fmt.Println("end!")
}
