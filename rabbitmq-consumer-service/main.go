package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	// 获取命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件的路径")
	flag.Parse()

	// 创建性能剖析文件
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Println("Could not create CPU profile:", err)
		return
	}
	defer cpuFile.Close()

	// 开启 CPU 性能剖析
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Println("Could not start CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	// 创建内存剖析文件
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Println("Could not create memory profile:", err)
		return
	}
	defer memFile.Close()

	// 定期生成内存剖析（每 30 秒一次）
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := pprof.WriteHeapProfile(memFile); err != nil {
					log.Println("Could not write memory profile:", err)
				} else {
					log.Println("Memory profile generated")
				}
			}
		}
	}()

	// 启动 HTTP 服务器并暴露 pprof 接口
	go func() {
		log.Println("Starting pprof server on :6060")
		err := http.ListenAndServe(":6060", nil) // 在 6060 端口提供 pprof 接口
		if err != nil {
			log.Fatalf("Could not start pprof server: %v", err)
		}
	}()

	// 加载配置
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %s", err)
	}

	// 创建消费者
	consumer, err := NewConsumer(config)
	if err != nil {
		log.Fatalf("创建消费者失败: %s", err)
	}
	defer consumer.Close()

	// 启动配置文件监控
	go func() {
		if err := watchConfigFile(*configPath, config); err != nil {
			log.Fatalf("监听配置文件时出错: %s", err)
		}
	}()

	// 启动心跳机制
	go sendHeartbeat(consumer.channel, config.RabbitMQ.QueueName)

	// 启动消费者
	log.Println("Starting consumers...")
	consumer.Start()
}
