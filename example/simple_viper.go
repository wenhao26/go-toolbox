package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("my")
	viper.SetConfigType("ini")
	viper.AddConfigPath("../conf/ini")

	/*err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败：%v", err))
	}*/
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
		} else {
			// 配置文件被找到，但产生了另外的错误
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", in.Name)
	})

	fmt.Println("before=", viper.Get("cent.addr"))
	time.Sleep(time.Second * 8)
	fmt.Println("after=", viper.Get("cent.addr"))
}
