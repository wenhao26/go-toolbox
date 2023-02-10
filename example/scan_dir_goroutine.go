package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var ch = make(chan struct{}, 255)

// 递归遍历
func walkDir(path string, fileSize chan<- int64) {
	defer wg.Done()

	// 限制并发量
	ch <- struct{}{}
	path = strings.Replace(path, "/", "\\", -1)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("读取文件夹出错：", err)
	}
	<-ch

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			log.Println("[DIR]=", fileInfo.Name())

			wg.Add(1)
			go walkDir(filepath.Join(path, fileInfo.Name()), fileSize)
		} else {
			log.Println("[FILENAME]=", path+fileInfo.Name())
			fileSize <- fileInfo.Size()
		}
	}
}

func main() {
	// 遍历所有文件目录，将其具体路径插入到MySQL中
	// 文件大小
	fileSize := make(chan int64)
	// 总文件大小
	var totalSize int64
	// 总文件数量
	var totalFiles int

	wg.Add(1)
	go walkDir("D:\\", fileSize)

	go func() {
		defer close(fileSize)
		wg.Wait()
	}()

	t := time.Now()
	for size := range fileSize {
		totalFiles++
		totalSize += size
	}

	fmt.Println("--目录文件总数：", totalFiles)
	// fmt.Println("--目录总大小:\t", totalSize/1024/1024, "M")
	fmt.Printf("--目录总大小:%.1fGB \n", float64(totalSize)/1e9)
	fmt.Println("--花费的时间：", time.Since(t).String())
}
