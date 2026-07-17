// Package config 负责解析命令行采参数，集中管理应用的可调优配置项
package config

import (
	"flag"
	"time"
)

// Config 整个生产者-消费者的运行参数
// 所有并发度、速率、缓冲区大小均可通过命令参数调整
// 便于在不同硬件环境下进行压测和调优，无需修改代码
type Config struct {
	ProducerNum     int           // 并发生产者数量，模拟多个实时数据源
	WorkerNum       int           // 消费者 Worker 数量，即消费者并发度
	QueueSize       int           // 生产者与消费者之间的缓冲队列容量（channel buffer）
	ProduceQPS      int           // 单个生产者每秒生产的消息数
	RunDuration     time.Duration // 程序运行时长，0 表示一直运行直到收到退出信号（Ctrl+C）
	ShutdownTimeout time.Duration // 优雅退出时，等待 worker 消费完剩余消息的最长时间
	StatsInterval   time.Duration // 统计信息打印周期
	LogLevel        string        // 日志级别：debug/info/warn/error
}

// Load 解析命令行参数并返回配置实例
func Load() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.ProducerNum, "producer", 3, "并发生产者数量")
	flag.IntVar(&cfg.WorkerNum, "worker", 8, "消费者 worker 数量")
	flag.IntVar(&cfg.QueueSize, "queue", 1000, "生产者-消费者之间的缓冲队列容量")
	flag.IntVar(&cfg.ProduceQPS, "qps", 200, "单个生产者每秒生产的消息数")
	flag.DurationVar(&cfg.RunDuration, "duration", 0, "运行时长，0 表示一直运行直到收到 Ctrl+C 信号")
	flag.DurationVar(&cfg.ShutdownTimeout, "shutdown-timeout", 10*time.Second, "优雅退出时等待 worker 消费完剩余消息的最长时间")
	flag.DurationVar(&cfg.StatsInterval, "stats-interval", 2*time.Second, "统计信息打印周期")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "日志级别：debug/info/warn/error")

	flag.Parse()

	return cfg
}
