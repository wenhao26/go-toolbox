package main

import (
	"log"
	"time"
)

func main() {
	// 初始化配置并开启监控
	InitConfig()
	log.Println("服务已启动，请尝试修改 config.yaml 中的配置项...")

	// 模拟常驻进程，不断读取配置
	for {
		currConf := GetConfig()

		log.Printf("[CurrentConfig] App: %v",
			currConf.App,
		)

		time.Sleep(5 * time.Second)
	}
}
