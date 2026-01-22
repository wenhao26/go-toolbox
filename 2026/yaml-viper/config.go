package main

import (
	"sync"
)

// GlobalConfig 全局配置单例
var (
	GlobalConfig = &Config{}
	configLock   = sync.RWMutex{} // 用于热更新时的并发安全
)

// Config 配置项结构体
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Logger LoggerConfig `mapstructure:"logger"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Jwt    JwtConfig    `mapstructure:"jwt"`
	Cron   CronConfig   `mapstructure:"cron"`
}

// AppConfig 应用配置项
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Prot    int    `mapstructure:"prot"`
	Version string `mapstructure:"version"`
}

// LoggerConfig 日志配置项
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// MySQLConfig MySQL配置项
type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置项
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JwtConfig JWT配置项
type JwtConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// CronConfig JWT配置项
type CronConfig struct {
	EnableSeconds bool `mapstructure:"enable_seconds"`
	Timeout       int  `mapstructure:"timeout"`
}

// GetConfig 安全地获取配置（防止热更新时读取到中间状态）
func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return GlobalConfig
}
