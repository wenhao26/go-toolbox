// Package stats 提供无锁的并发安全统计计数器，用于实时观测系统吞吐与健康状况
package stats

import (
	"sync/atomic"
)

// Stats 使用 atomic.Uint64 而非 mutex + int64，避免在高并发场景下
// 因加锁产生的竞争开销，是 Go 中“高频写、低频读”计数器的标准实践
type Stats struct {
	Produced  atomic.Uint64 // 累计生产消息数
	Consumed  atomic.Uint64 // 累计消费消息数
	Succeeded atomic.Uint64 // 累计处理成功数
	Failed    atomic.Uint64 // 累计处理失败/被中断数（含 panic 恢复、超时兜底）
}

// New 创建一个初始为零值的 Stats
func New() *Stats {
	return &Stats{}
}

// IncProduced 生产计数 +1
func (s *Stats) IncProduced() { s.Produced.Add(1) }

// IncConsumed 消费计数 +1
func (s *Stats) IncConsumed() { s.Consumed.Add(1) }

// IncSucceeded 处理成功计数 +1
func (s *Stats) IncSucceeded() { s.Succeeded.Add(1) }

// IncFailed 处理失败计数 +1
func (s *Stats) IncFailed() { s.Failed.Add(1) }

// Snapshot 原子读取当前所有计数器的快照值，用于打印或者上报监控系统
func (s *Stats) Snapshot() (produced, consumed, Succeeded, failed uint64) {
	return s.Produced.Load(), s.Consumed.Load(), s.Succeeded.Load(), s.Failed.Load()
}
