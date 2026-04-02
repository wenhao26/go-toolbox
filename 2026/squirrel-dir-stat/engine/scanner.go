package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// -- 核心扫描引擎 --

// Scanner 扫描器
type Scanner struct {
	Counter *AtomicCounter
	wg      sync.WaitGroup
	jobChan chan string // 任务通道：传递需要扫描的目录路径
}

// NewScanner 实例化扫描器
func NewScanner(workerCount int) *Scanner {
	return &Scanner{
		Counter: &AtomicCounter{},
		jobChan: make(chan string, workerCount), // 足够大的缓冲区，防止初期阻塞
	}
}

// Scan 执行扫描任务
func (s *Scanner) Scan(root string) {
	// 启动固定数量的工作写成
	for i := 0; i < 64; i++ { // TODO：生产环境下，建议根据 CPU 核数和磁盘 IO 设置
		go s.worker()
	}

	// 提交初始根目录任务
	s.submit(root)

	// 等待所有层级的目录扫描完成
	s.wg.Wait()
	close(s.jobChan)
}

// worker 扫描工作者
func (s *Scanner) worker() {
	for path := range s.jobChan {
		s.processDir(path)
		s.wg.Done()
	}
}

// processDir 处理目录
func (s *Scanner) processDir(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// 容错处理：记录但跳过无权限或特殊系统文件
		fmt.Printf("[Error] 无法访问目录 %s: %v\n", dirPath, err)
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// 如果是目录，递归提交到任务池
			s.submit(fullPath)
		} else {
			// 如果是文件，记录路径并原子统计
			info, err := entry.Info()
			if err != nil {
				continue
			}

			s.Counter.Add(info.Size())

			// TODO：实时输出扫描路径（高并发下，此操作会稍微降低性能，生成环境可选择性开启即可）
			fmt.Println("Scanning:", fullPath)
		}
	}
}

// submit 提交属于根目录到任务池
func (s *Scanner) submit(path string) {
	s.wg.Add(1)
	s.jobChan <- path
}
