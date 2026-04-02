package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"toolbox/2026/squirrel-dir-stat/engine"
)

func main() {
	pathPtr := flag.String("p", "", "待扫描的目录路径")
	workerPtr := flag.Int("w", runtime.NumCPU()*8, "并发工作协程数")
	flag.Parse()

	if *pathPtr == "" {
		fmt.Println("用法: squirrel-stat -p <目录路径> [-w <并发数>]")
		os.Exit(1)
	}

	// 检查目录有效性
	info, err := os.Stat(*pathPtr)
	if err != nil || !info.IsDir() {
		log.Fatalf("错误: 路径 '%s' 不是有效的目录\n", *pathPtr)
	}

	fmt.Printf("🐿️ 统计松鼠启动! 目标: %s, 并发权重: %d\n", *pathPtr, *workerPtr)
	fmt.Println("--------------------------------------------------")

	startTime := time.Now()

	// 初始化并运行引擎
	scanner := engine.NewScanner(*workerPtr)
	scanner.Scan(*pathPtr)

	// 结果输出
	elapsed := time.Since(startTime)
	res := scanner.Counter.Snapshot()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("✅ 统计完成!\n")
	fmt.Printf("⏳ 总耗时: %v\n", elapsed)
	fmt.Printf("📂 总文件数: %d\n", res.TotalFileCount)
	fmt.Printf("💾 总大小: %s\n", engine.FormatSize(res.TotalFileSize))
}
