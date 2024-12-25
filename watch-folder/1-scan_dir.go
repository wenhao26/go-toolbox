package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

var query = "config.php"
var matches int
var workerCount = 0                   // 工作者数
var maxWorkerCount = 15               // 最大工作者数
var searchRequest = make(chan string) // 查询请求管道
var workerDone = make(chan bool)      // 工作者完成任务管道
var foundMatch = make(chan bool)      // 发现匹配管道

func search(path string, isMaster bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()
		if name == query {
			foundMatch <- true
		}
		if file.IsDir() {
			if workerCount <= maxWorkerCount {
				searchRequest <- path + name + "\\"
			} else {
				search(path+name+"\\", false)
			}
		}
	}

	if isMaster {
		workerDone <- true
	}
}

func waitForWorker() {
	for {
		select {
		case path := <-searchRequest:
			workerCount++
			go search(path, true)
		case <-workerDone:
			workerCount--
			if workerCount == 0 {
				return
			}
		case <-foundMatch:
			matches++
		}
	}
}

func main() {
	startTime := time.Now()
	go search("D:\\www\\", true)
	waitForWorker()
	fmt.Println(matches, " matches.")
	fmt.Println(time.Since(startTime))
}
