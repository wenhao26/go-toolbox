package main

import (
	"flag"
	"fmt"
	"time"

	"toolbox/conf"
)

// 心知天气

func main() {
	var city string
	flag.StringVar(&city, "city", "", "查询某个城市名称拼音，如：./main --city=guangzhou")
	flag.Parse()

	cfg := conf.GetINI()
	section := cfg.Section("xinzhi_weather")

	//
	apiURL := fmt.Sprintf("https://api.seniverse.com/v3/weather/daily.json?key=%s&location=%s&language=zh-Hans&unit=c&start=-1&days=3",
		section.Key("api_key").String(),
		city,
	)
	fmt.Println(apiURL)

	// 输出不换行，同一行刷新显示信息
	var index = 0
	for {
		index++
		fmt.Printf("\rhalo %d%%", index)
		time.Sleep(1e9)
	}
}
