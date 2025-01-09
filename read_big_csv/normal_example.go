package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func readCSVFile() {
	// 打开CSV文件
	file, err := os.Open("F:\\test_files\\csv\\pp-complete.csv") // 4.84GB
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	startTime := time.Now()
	recordCount := 0

	for {
		record, err := reader.Read() // 读取一行
		if err != nil {
			break
		}

		// 打印每一行的内容
		// 为了避免输出过多，选择限制打印的行数
		if recordCount%1000 == 0 { // 每间隔1000行打印一次
			fmt.Println(record)
		}

		// 处理每一行数据
		recordCount++
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("普通方式读取CSV文件，共处理了 %d 行数据，耗时 %s\n", recordCount, elapsedTime)
}

func main() {
	readCSVFile()
}
