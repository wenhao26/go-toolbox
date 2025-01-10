package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 配置结构体
type Config struct {
	RabbitMQ struct {
		URL            string `yaml:"url"`
		QueueName      string `yaml:"queue_name"`
		PrefetchCount  int    `yaml:"prefetch_count"`
		ConsumerCount  int    `yaml:"consumer_count"`
		WorkerPoolSize int    `yaml:"worker_pool_size"`
	} `yaml:"rabbitmq"`
}

// 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
