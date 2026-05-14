package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// FileInfo 文件信息
type FileInfo struct {
	Path string
	Size int64
}

func main() {
	paths := make(chan string, 100)
	results := make(chan FileInfo, 100)

	// 开启协程，统计结果
	var totalSize int64
	var fileCount int
	done := make(chan struct{})
	go func() {
		for file := range results {
			totalSize += file.Size
			fileCount++
			// fmt.Printf("已统计: %s (%d bytes)\n", file.Path, file.Size)
		}
		done <- struct{}{} // 统计完，发信号给主线程
	}()

	// 启动工作池
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for path := range paths {
				// TODO 模拟扫描磁盘的耗时和结果
				time.Sleep(500 * time.Millisecond)
				results <- FileInfo{Path: path, Size: 1024 * 1024}

				fmt.Printf("[Worker%d] 完成扫描: %s\n", workerID, path)
			}
		}(i)
	}

	// TODO 模拟 20 个文件夹需要统计
	go func() {
		for i := 0; i < 20; i++ {
			paths <- fmt.Sprintf("/data/project/dir_%d", i)
		}
		close(paths) // 分派完成，关闭任务通道
	}()

	wg.Wait()
	close(results)

	<-done
	fmt.Printf("📊 统计报告：\n文件总数: %d\n总大小: %.2f MB\n", fileCount, float64(totalSize)/1024/1024)
}
