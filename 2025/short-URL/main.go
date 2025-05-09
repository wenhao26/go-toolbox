package main

import (
	"fmt"

	"toolbox/2025/short-URL/api"
)

// -- 示例 --
// 创建短链：curl -X POST http://localhost:8080/shorten -d '{"long_url":"待生成的长链接"}'

func main() {
	fmt.Println("短链URL服务启动...")
	api.StartServer()
}
