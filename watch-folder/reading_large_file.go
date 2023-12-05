package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Running...")
	//filename := "F:\\test_files\\test.txt"
	filename := "F:\\test_files\\novel.txt"
	//filename := "F:\\test_files\\729mb.txt"
	chunkRead(filename)

	fmt.Println("chunkRead Spend:", time.Now().Sub(start))
}

func chunkRead(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()

	// 设置每次读取的字节数
	buffer := make([]byte, 10*1024*1024)
	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		if n == 0 {
			break
		}

		// todo...
		//fmt.Println(string(buffer[:n]))
	}
}
