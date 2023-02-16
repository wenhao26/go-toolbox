package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"

	"toolbox/conf"
)

// 心知天气数据查询
func weather(apiKey, city string) {
	url := fmt.Sprintf("https://api.seniverse.com/v3/weather/daily.json?key=%s&location=%s&language=zh-Hans&unit=c&start=-1&days=5",
		apiKey,
		city,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonMap := make(map[string]interface{})
	_ = json.Unmarshal(body, &jsonMap)

	data := [][]string{}
	if v, ok := jsonMap["results"].([]interface{})[0].(map[string]interface{}); ok {
		for _, daily := range v["daily"].([]interface{}) {
			dailyItem := daily.(map[string]interface{})
			data = append(data, []string{
				dailyItem["date"].(string),
				dailyItem["text_day"].(string),
				dailyItem["text_night"].(string),
				dailyItem["high"].(string),
				dailyItem["low"].(string),
				dailyItem["wind_direction"].(string),
				dailyItem["wind_direction_degree"].(string),
				dailyItem["wind_speed"].(string),
				dailyItem["wind_scale"].(string),
				dailyItem["humidity"].(string),
			})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"日期", "白天", "夜晚", "最高温", "最低温", "风向", "风向度", "风速", "风力刻度", "湿度"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func
main() {
	var city string
	flag.StringVar(&city, "city", "", "查询某个城市名称，如：./main --city=guangzhou")
	flag.Parse()

	cfg := conf.GetINI()
	section := cfg.Section("xinzhi_weather")
	apiKey := section.Key("api_key").String()

	// weather(apiKey, city)

	// 输出不换行，同一行刷新显示信息
	index := 0
	for {
		index++
		select {
		case <-time.After(60 * time.Second):
			//fmt.Printf("\r >>> 心知天气数据已被调用%d次 <<<", index)
			log.Printf("心知天气数据提醒您，查询城市-[%s]天气情况已被刷新！", city)
			weather(apiKey, city)
			fmt.Println("")
		}
	}
}
