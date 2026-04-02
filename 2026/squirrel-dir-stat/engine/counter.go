package engine

import (
	"fmt"
	"sync/atomic"
)

// -- 高性能统计器 --

// Result 存储最终统计结果
type Result struct {
	TotalFileCount int64
	TotalFileSize  int64
}

// AtomicCounter 使用原子操作保证并发下的高性能统计
type AtomicCounter struct {
	fileCount int64
	fileSize  int64
}

// Add 累加文件和大小
func (c *AtomicCounter) Add(size int64) {
	atomic.AddInt64(&c.fileCount, 1)
	atomic.AddInt64(&c.fileSize, size)
}

// Snapshot 获取当前数值快照
func (c *AtomicCounter) Snapshot() Result {
	return Result{
		TotalFileCount: atomic.LoadInt64(&c.fileCount),
		TotalFileSize:  atomic.LoadInt64(&c.fileSize),
	}
}

// FormatSize 转换字节格式
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
