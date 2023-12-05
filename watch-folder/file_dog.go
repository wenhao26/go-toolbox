package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	wg       sync.WaitGroup
	sum      = 0
	filename = "F:\\test_files\\test.txt"
	//filename = "F:\\test_files\\novel.txt"
	//filename = "F:\\test_files\\small.txt"
)

func main() {
	// 一次性将数据写入内存处理方式
	/*b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("read failure:", err)
	}

	// 逐行打印
	text := string(b)
	lineList := strings.Split(text, "\n\r")
	for _, str := range lineList {
		fmt.Println(str)
	}*/

	// 基于磁盘和缓存处理方式
	/*f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("read failure:", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Println(string(line))
	}*/

	// 统计文本字数
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("read failure:", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		fmt.Println(string(line))
		wg.Add(1)
		go sumWC(string(line))
	}
	wg.Wait()

	fmt.Println("\n\r>> Number of words:", sum)
}

func sumWC(str string) {
	sum += len(str)
	wg.Done()
}
