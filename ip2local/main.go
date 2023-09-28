package main

import (
	"flag"
	"fmt"

	"toolbox/ip2local/ip"
)

func main() {
	var action string
	var ipAddress string

	flag.StringVar(&action, "a", "", "操作动作:pull(拉取)，create(创建)，save(保存)，search(搜索)")
	flag.StringVar(&ipAddress, "ip", "", "查询的IP地址")
	flag.Parse()

	server := ip.New()

	switch action {
	case "pull":
		server.Pull()
	case "create":
		server.Create()
	case "save":
		server.SaveDb()
	case "search":
		if ipAddress == "" {
			fmt.Println("请输入 -ip {你要查询的ip}")
			return
		}
		server.Search(ipAddress)
	default:
		fmt.Println("执行 -h 命令查看")
	}
}
