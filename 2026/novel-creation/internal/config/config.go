package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Redis  RedisConfig
	Queue  QueueConfig
	Worker WorkerConfig
	Log    LogConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

type QueueConfig struct {
	Name             string
	ProcessingSuffix string
	BlockTimeout     int // 秒
}

type WorkerConfig struct {
	Count           int
	GracefulTimeout int // 秒
}

type LogConfig struct {
	Level  string
	Output string
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
			PoolSize: viper.GetInt("redis.pool_size"),
		},
		Queue: QueueConfig{
			Name:             viper.GetString("queue.name"),
			ProcessingSuffix: viper.GetString("queue.processing_suffix"),
			BlockTimeout:     viper.GetInt("queue.block_timeout"),
		},
		Worker: WorkerConfig{
			Count:           viper.GetInt("worker.count"),
			GracefulTimeout: viper.GetInt("worker.graceful_timeout"),
		},
		Log: LogConfig{
			Level:  viper.GetString("log.level"),
			Output: viper.GetString("log.output"),
		},
	}

	// 设置默认值
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 10
	}
	if cfg.Queue.BlockTimeout == 0 {
		cfg.Queue.BlockTimeout = 5
	}
	if cfg.Worker.Count == 0 {
		cfg.Worker.Count = 5
	}
	if cfg.Worker.GracefulTimeout == 0 {
		cfg.Worker.GracefulTimeout = 30
	}
	return cfg, nil
}
