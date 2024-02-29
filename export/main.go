package main

import (
	"fmt"
	"time"

	"toolbox/export/getdata"
)

func main() {
	start := time.Now().UnixNano() / 1e6
	getdata.Export()
	end := time.Now().UnixNano() / 1e6
	// fmt.Printf("测试--用时:%d毫秒\r\n", end-start)
	fmt.Printf("测试--用时:%d秒\r\n", (end-start)/1000)
}
