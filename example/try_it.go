package main

import (
	"flag"
	"fmt"
	"time"
)

// 输出时间范围区间内的日期列表
func rangeIntervalDate(sDate, eDate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sDate) {
		timeFormatTpl = timeFormatTpl[0:len(sDate)]
	}
	date, err := time.Parse(timeFormatTpl, sDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, eDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

func main() {
	var optName string
	flag.StringVar(&optName, "opt", "", "输入对应业务场景的操作名称进行执行，如：./main --opt=test")
	flag.Parse()

	if optName == "rangeIntervalDate" {
		fmt.Println("说明：输出时间范围区间内的日期列表")
		ret := rangeIntervalDate("2023-01-25", "2023-02-14")
		fmt.Println(ret)
	}

}
