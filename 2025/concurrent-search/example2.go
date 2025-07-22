package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

// - 一个高性能的并发文件检索工具的实现，支持模糊匹配、优雅退出和大目录高效检索（DeepSeek）
// - /fileSearch -dir /path/to/search -pattern "*.go"

// FileSearcher 文件搜索器的相关参数和状态
type FileSearcher struct {
	rootDir    string         // 要搜索的根目录
	pattern    string         // 要匹配的文件名（支持模糊匹配）
	concurrent int            // 并发协程数量
	results    chan string    // 用于收集结果的通道
	done       chan struct{}  // 用于通知停止搜索的通道
	sigChan    chan os.Signal // 系统信号通道
	wg         sync.WaitGroup // 用于等待所有协程完成
	mu         sync.Mutex     // 用于保护共享状态
	foundCount int            // 已找到的文件计数
}

// search 启动并发文件搜索
func (fs *FileSearcher) search() {
	guard := make(chan struct{}, fs.concurrent) // 限制并发目录读取的协程数量

	// 启动第一个协程处理根目录
	fs.wg.Add(1)
	go fs.walkDir(fs.rootDir, guard)

	// 等待所有工作完成
	go func() {
		fs.wg.Wait()
		close(fs.results)
	}()
}

// walkDir 递归遍历目录
func (fs *FileSearcher) walkDir(dir string, guard chan struct{}) {
	defer fs.wg.Done()

	select {
	case <-fs.done: // 检查是否应该停止
		return
	default:
	}

	// 获取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// 处理目录中的每个条目
	for _, entry := range entries {
		select {
		case <-fs.done:
			return
		default:
		}

		fullPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// 如果是目录，启动新的协程处理它（如果未达到并发限制）
			select {
			case guard <- struct{}{}: // 尝试获取一个槽位
				fs.wg.Add(1)
				go func(path string) {
					defer func() { <-guard }() // 完成后释放槽位
					fs.walkDir(path, guard)
				}(fullPath)
			default:
				// 如果并发限制一达到，则在当前协程中处理
				fs.wg.Add(1)
				fs.walkDir(fullPath, guard)
			}
		} else {
			// 如果是文件，检查是否匹配
			if fs.matchPattern(entry.Name()) {
				fs.results <- fullPath
			}
		}
	}
}

// matchPattern 检查文件名是否匹配搜索模式
func (fs *FileSearcher) matchPattern(filename string) bool {
	pattern := fs.pattern

	// 如果文件名为`*`或为空，匹配所有文件
	if pattern == "*" || pattern == "" {
		return true
	}

	if strings.Contains(pattern, "*") {
		// 将模式转换为正则表达式
		pattern = strings.ReplaceAll(pattern, "*", ".*")
		return strings.Contains(filename, strings.ReplaceAll(pattern, ".*", ""))
	}

	// 包含匹配
	return strings.Contains(filename, pattern)
}

// processResults 处理搜索结果并监听退出信号
func (fs *FileSearcher) processResults() {

	for {
		select {
		case sig := <-fs.sigChan:
			// 收到中断信号，优雅退出
			fmt.Printf("\n收到信号 %v, 正在停止搜索...\n", sig)
			close(fs.done) // 通知所有协程停止
			return
		case result, ok := <-fs.results:
			if !ok {
				return
			}
			fs.mu.Lock()
			fs.foundCount++
			fs.mu.Unlock()

			fmt.Println("- ", result)
		}
	}
}

func main() {
	// 解析命名行参数
	dir := flag.String("dir", "", "要搜索的目录")
	pattern := flag.String("filename", "", "要搜索的文件名（支持模糊匹配）")
	workers := flag.Int("workers", runtime.NumCPU()*2, "并发工作协程数量")
	flag.Parse()

	if *dir == "" {
		fmt.Println("错误：必须指定搜索的目录")
		flag.Usage()
		os.Exit(1)
	}
	if *pattern == "" {
		fmt.Println("错误：必须指定搜索的文件")
		flag.Usage()
		os.Exit(1)
	}

	// 初始化文件搜索器
	searcher := &FileSearcher{
		rootDir:    *dir,
		pattern:    *pattern,
		concurrent: *workers,
		results:    make(chan string, 1000),
		done:       make(chan struct{}),
		sigChan:    make(chan os.Signal, 1),
	}

	// 设置信号处理
	signal.Notify(searcher.sigChan, syscall.SIGINT, syscall.SIGTERM)

	start := time.Now()

	// 启动搜索
	go searcher.search()

	// 处理结果和信号
	searcher.processResults()

	duration := time.Since(start)
	fmt.Printf("\n搜索完成！共找到 %d 个文件, 耗时 %v\n", searcher.foundCount, duration)
}
