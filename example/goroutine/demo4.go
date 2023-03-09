package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func uploadFile(uploadResultCh chan<- string) {
	go func() {
		for {
			filename := genFilename()
			fmt.Printf("[ready]正在上传文件[%s]到服务器...\n", filename)
			uploadResultCh <- fmt.Sprintf("[ok]文件[%s]已存储成功！", filename)
		}
	}()
}

func pullUploadResult(uploadResultCh <-chan string) {
	go func() {
		for {
			select {
			case val, ok := <-uploadResultCh:
				if !ok {
					return
				}
				log.Println(val)
			}
		}
	}()
}

func genFilename() string {
	t := time.Now().Unix()
	ext := ".jpg"
	return strconv.Itoa(int(t)) + ext
}

func main() {
	uploadResultCh := make(chan string, 100)
	uploadFile(uploadResultCh)
	pullUploadResult(uploadResultCh)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
