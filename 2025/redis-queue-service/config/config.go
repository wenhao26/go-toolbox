package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// RedisConfig 连接Redis所需的配置信息
type RedisConfig struct {
	Addr         string        // 服务器地址
	Password     string        // 连接密码
	DB           int           // 数据库索引值
	PoolSize     int           // 连接池最大连接数
	DialTimeout  time.Duration // 连接超时时间
	ReadTimeout  time.Duration // 读取超时时间
	WriteTimeout time.Duration // 写入超时时间
}

// QueueConfig 队列服务配置信息
type QueueConfig struct {
	Name         string        // Redis列表用于队列键名
	Concurrency  int           // 消费者并发处理消息的协程数量
	BlockTimeout time.Duration // BRPOP命令的阻塞超时时间，0表示无阻塞
}

// Config 服务配置集合
type Config struct {
	Redis RedisConfig
	Queue QueueConfig
}

// LoadConfig 从环境变量加载配置
// 推荐使用该方式，方便在不同环境部署
func LoadConfig() (*Config, error) {
	// 加载 .env 文件中的环境变量
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// 加载Redis配置
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0
	}

	redisPoolSizeStr := os.Getenv("REDIS_POOL_SIZE")
	redisPoolSize, err := strconv.Atoi(redisPoolSizeStr)
	if err != nil || redisPoolSize == 0 {
		redisPoolSize = 10
	}

	// 加载队列配置
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "my_redis_queue"
	}

	concurrencyStr := os.Getenv("CONSUMER_CONCURRENCY")
	concurrency, err := strconv.Atoi(concurrencyStr)
	if err != nil || concurrency == 0 {
		concurrency = 5
	}

	blockTimeoutSecStr := os.Getenv("QUEUE_BLOCK_TIMEOUT_SECONDS")
	blockTimeoutSec, err := strconv.Atoi(blockTimeoutSecStr)
	if err != nil {
		blockTimeoutSec = 5
	}

	return &Config{
		Redis: RedisConfig{
			Addr:         redisAddr,
			Password:     redisPassword,
			DB:           redisDB,
			PoolSize:     redisPoolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		Queue: QueueConfig{
			Name:         queueName,
			Concurrency:  concurrency,
			BlockTimeout: time.Duration(blockTimeoutSec) * time.Second,
		},
	}, nil
}

// ValidateConfig 验证配置有效性
func (c *Config) ValidateConfig() error {
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis地址不能为空")
	}
	if c.Queue.Name == "" {
		return fmt.Errorf("队列名称不能为空")
	}
	if c.Queue.Concurrency <= 0 {
		return fmt.Errorf("消费者并发值必须大于0")
	}
	if c.Queue.BlockTimeout < 0 {
		return fmt.Errorf("队列块超时不能为负")
	}
	return nil
}
