package main

import (
	"fmt"
	"sync"
	"time"
)

// - 实际应用：文件批量处理器

// FileProcessor 文件处理器
type FileProcessor struct {
	done  chan struct{}
	wg    sync.WaitGroup
	files chan string
}

// NewFileProcessor 创建一个新的文件处理器
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		done:  make(chan struct{}),
		files: make(chan string, 100),
	}
}

// Start 启动处理器
func (fp *FileProcessor) Start(workers int) {
	for i := 0; i < workers; i++ {
		fp.wg.Add(1)
		go fp.worker(i)
	}
}

// worker 工作者
func (fp *FileProcessor) worker(id int) {
	defer fp.wg.Done()

	for {
		select {
		case <-fp.done:
			fmt.Printf("Worker %d 正在退出\n", id)
			return
		case file := <-fp.files:
			fmt.Printf("Worker %d 正在处理文件: %s\n", id, file)
			time.Sleep(200 * time.Millisecond) // 模拟处理时间
		}
	}
}

// AddFile 添加文件
func (fp *FileProcessor) AddFile(filename string) {
	select {
	case fp.files <- filename:
	case <-fp.done:
	}
}

// Stop 停止工作者
func (fp *FileProcessor) Stop() {
	close(fp.done)
	fp.wg.Wait()
	fmt.Println("所有 worker 已停止")
}

func main() {
	processor := NewFileProcessor()
	processor.Start(3)

	// 添加文件
	for i := 0; i < 10; i++ {
		processor.AddFile(fmt.Sprintf("file_%d.txt", i))
	}

	time.Sleep(1 * time.Second)
	processor.Stop()
}
