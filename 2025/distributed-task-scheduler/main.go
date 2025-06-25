package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"toolbox/2025/distributed-task-scheduler/master"
	"toolbox/2025/distributed-task-scheduler/worker"
)

func main() {
	// 定义一个命令行参数mode，用于指定节点类型（主节点或工作节点）
	var mode string
	flag.StringVar(&mode, "mode", "master", "mode of the node (master or worker)")
	flag.Parse()

	var err error

	// 根据mode参数值，启动主节点或者工作节点
	if mode == "master" {
		err = master.StartMaster()
	} else if mode == "worker" {
		err = worker.StartWorker()
	} else {
		log.Fatalf("invalid mode:%s", mode)
	}

	if err != nil {
		log.Fatalf("failed to start %s node: %v", mode, err) // 如果启动失败，记录错误并退出
	}

	// 设置优雅关闭信号处理
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop // 等待接收大盘终止信号

	// 根据mode参数，停止主节点或者工作节点
	if mode == "master" {
		master.StopMaster()
	} else {
		worker.StopWorker()
	}
}
