package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// mu 用于保护 timer 的并发访问
	mu    sync.Mutex
	timer *time.Timer
)

const delay = 100 * time.Millisecond // 防抖延迟时间

// InitConfig 初始化配置
func InitConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// 首次读取
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置失败: %s", err))
	}

	// 解析到全局变量
	unmarshalConfig(v)

	// 动态监控配置文件
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("检测到配置文件变更: %s, 操作类型: %s", e.Name, e.Op)

		mu.Lock()
		defer mu.Unlock()

		// 如果 100ms 内再次触发，停止之前的计时器
		if timer != nil {
			timer.Stop()
		}

		// 重新开始计时
		timer = time.AfterFunc(delay, func() {
			reload(v, e)
		})
	})
}

// unmarshalConfig 解析配置
func unmarshalConfig(v *viper.Viper) {
	if err := v.Unmarshal(GlobalConfig); err != nil {
		panic(fmt.Errorf("初始解析配置失败: %s", err))
	}
}

// reload 真正执行配置重载的操作
func reload(v *viper.Viper, e fsnotify.Event) {
	// 加锁保护并重新解析
	configLock.Lock()
	defer configLock.Unlock()

	// 避坑关键点：重载前务必重新 ReadInConfig，否则内存中的缓存可能没更新
	if err := v.ReadInConfig(); err != nil {
		log.Printf("热更新读取失败: %v", err)
		return
	}

	if err := v.Unmarshal(GlobalConfig); err != nil {
		log.Printf("热更新解析失败: %v", err)
	} else {
		log.Printf("[HotReload] 配置已自动更新。详情: %v", GlobalConfig)
	}
}
