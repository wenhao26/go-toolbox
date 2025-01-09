package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func readCSVFileConcurrently(filePath string, startLine, endLine int, wg *sync.WaitGroup, resultChan chan int) {
	defer wg.Done()

	// 打开CSV文件（仅打开一次）
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 跳过前startLine行
	for i := 0; i < startLine; i++ {
		_, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
	}

	recordCount := 0
	for {
		// 读取一行
		record, err := reader.Read()
		if err != nil {
			break
		}

		// 如果已经读取到endLine行，停止读取
		if recordCount+startLine >= endLine {
			break
		}

		// 打印每一行的内容
		// 为了避免输出过多，选择限制打印的行数
		if recordCount%1000 == 0 { // 每间隔1000行打印一次
			fmt.Println(record)
		}

		// 统计行数
		recordCount++
	}

	// 将结果发送到channel
	resultChan <- recordCount
}

func readCSVFile2() {
	// 文件路径
	filePath := "F:\\test_files\\csv\\pp-complete.csv"

	// 打开CSV文件获取文件信息
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 获取文件总行数
	var totalLines int
	reader := csv.NewReader(file)
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}
		totalLines++
	}

	// 假设将文件分成 2 或 4 个部分，避免协程过多
	numChunks := 2
	linesPerChunk := totalLines / numChunks

	startTime := time.Now()

	var wg sync.WaitGroup
	resultChan := make(chan int, numChunks)

	for i := 0; i < numChunks; i++ {
		// 每个协程处理的起始行和结束行
		startLine := i * linesPerChunk
		endLine := (i + 1) * linesPerChunk
		if i == numChunks-1 {
			// 最后一部分处理到文件末尾
			endLine = totalLines
		}

		wg.Add(1)
		go readCSVFileConcurrently(filePath, startLine, endLine, &wg, resultChan)
	}

	// 等待所有协程完成
	wg.Wait()

	// 聚合结果
	totalRecords := 0
	close(resultChan)
	for count := range resultChan {
		totalRecords += count
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("并发方式读取CSV文件，共处理了 %d 行数据，耗时 %s\n", totalRecords, elapsedTime)
}

func main() {
	readCSVFile2()
}
