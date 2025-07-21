package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// - 高性能文件搜索器，支持在指定目录下进行模糊检索
// - 使用 go 的并发特性，通过扇出（Fan-out）/扇入（Fan-in）模式来并行检索
// - 版本v1：潜在的瓶颈是目录遍历本身是单线程的，虽然工作协程可以并行检索文件，但只有一个协程在分发目录任务。
//			对于包含大量子目录的巨大文件系统来说，这个分发者可能成为瓶颈。

// searchFilesConcurrentlyV1 并发模型查找匹配文件
func searchFilesConcurrentlyV1(rootPath, filename string) ([]string, error) {
	jobs := make(chan string, 100)    // 用于分发目录给工作协程
	results := make(chan string, 100) // 用于收集找到的文件路径

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 工作协程循环，从任务信道中读取目录
			for dirPath := range jobs {
				entries, err := os.ReadDir(dirPath)
				if err != nil {
					continue
				}

				for _, entry := range entries {
					fmt.Println(entry.Name())
					// 使用 filepath.Match 进行模糊匹配
					match, err := filepath.Match(filename, entry.Name())
					if err != nil {
						continue
					}
					if match {
						// 匹配成功，将完整的路径发送为结果信道
						results <- filepath.Join(dirPath, entry.Name())
					}
				}
			}
		}()
	}

	// 启动一个独立协程来遍历目录树（类似生产者）
	go func() {
		// 使用 filepath.WalkDir 遍历目录树，将每个子目录作为任务放入 jobs 信道
		_ = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				jobs <- path
			}
			return nil
		})
		close(jobs) // 目录遍历完后，关闭 jobs 信道，通知工作协程没有更多任务
	}()

	// 启动一个独立协程来等待所有工作协程完成，然后关闭 results 信道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 主协程从 results 信道中收集所有的结果（扇入）
	var foundFiles []string
	for file := range results {
		foundFiles = append(foundFiles, file)
	}

	return foundFiles, nil
}

// 使用：go run main.go -dir {检索目录路径} -filename {检索文件}
func main() {
	var dir string
	var filename string

	flag.StringVar(&dir, "dir", "", "检索目录")
	flag.StringVar(&filename, "filename", "", "检索文件名，支持模糊查询，如：*.txt")
	flag.Parse()

	if filename == "" {
		fmt.Println("必须提供一个要搜索的文件名")
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()

	files, err := searchFilesConcurrentlyV1(dir, filename)
	if err != nil {
		fmt.Printf("搜索失败：%v\n", err)
		return
	}

	duration := time.Since(start)
	fmt.Printf("\n结果：在目录[%s]下找到[%d]个匹配文件，用时[%v]\n", dir, len(files), duration)

	for _, file := range files {
		fmt.Println(file)
	}
}
