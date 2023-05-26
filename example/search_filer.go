package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 搜索文件器 - v1.0
func main() {
	var path string
	flag.StringVar(&path, "path", "", "待扫描的目录地址。如：./main --path=/data/www")
	flag.Parse()

	t := time.Now()

	if path == "" {
		log.Fatalln("请输入目录地址")
	}

	filename := ""
	_, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}

	var fileCount int
	fileCountCh := make(chan int)

	// 开启协程扫描目录
	var wg sync.WaitGroup
	findFile(path, filename, fileCountCh, &wg)
	go func() {
		defer close(fileCountCh)
		wg.Wait()
	}()

	for n := range fileCountCh {
		fileCount += n
	}
	fmt.Println("--共扫描文件：", fileCount)
	fmt.Println("--消耗时间：", time.Since(t).String())

	/*c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c*/
}

func findFile(path, filename string, fileCountCh chan int, wg *sync.WaitGroup) {
	path = strings.Replace(path, "/", "\\", -1)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("读取文件夹出错：", err)
	}

	for _, fileInfo := range fileInfos {
		wg.Add(1)
		go func(fileInfo fs.FileInfo) {
			defer wg.Done()
			if fileInfo.IsDir() {
				dir := filepath.Join(path, fileInfo.Name())
				log.Println("[DIR]=", dir)
				findFile(dir, filename, fileCountCh, wg)
			} else {
				fileCountCh <- 1
				log.Println(" --FILENAME:", filepath.Join(path, fileInfo.Name()))
			}
		}(fileInfo)
	}
}
