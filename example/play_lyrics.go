package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func readLRC(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseLRC(text string) (endMs int, lyric string) {
	// 分离格式标签和歌词
	result := strings.Split(text, "]")
	format, lyric := result[0], result[1]

	// 去除格式标签左首位[
	format = strings.TrimLeft(format, "[")

	// 计算歌词播报的毫秒时间
	var min, second, ms, totalMs int
	formatArr := strings.Split(format, ":")
	lTag, rTag := formatArr[0], formatArr[1]
	if lTag != "" || rTag != "" {
		min, _ = strconv.Atoi(lTag)

		// 判断右侧标签是否存在毫秒情况
		if ok := strings.Contains(rTag, "."); ok {
			rTagTmp := strings.Split(rTag, ".")
			second, _ = strconv.Atoi(rTagTmp[0])
			ms, _ = strconv.Atoi(rTagTmp[1])
		} else {
			second, _ = strconv.Atoi(rTag)
		}
		totalMs = min*60*1000 + second*1000 + ms
	}
	return totalMs, lyric
}

func main() {
	var lastMs int

	// 加载歌词
	data, err := readLRC("data/lrc/piaoxiangbeifang.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// 解析歌词
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); i++ {
		endMs, lyric := parseLRC(lines[i])
		time.Sleep(time.Duration(endMs-lastMs) * time.Millisecond)
		lastMs = endMs

		if len(lyric) > 1 {
			fmt.Println("  ♪", lyric)
		}
	}

	fmt.Println(" >>> 歌词播放结束")
}
