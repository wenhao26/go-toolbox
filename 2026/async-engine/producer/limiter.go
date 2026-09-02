package producer

import (
	"sync/atomic"
	"time"
)

// Limiter 基于原子令牌桶实现非阻塞速率限制，支持运行时动态调整目标速率。
// 所有操作均为无锁实现，可在高并发场景下安全使用。
type Limiter struct {
	rate   int64 // 目标速率（token/s）
	tokens int64 // 当前可用令牌数
	last   int64 // 上次补充时间（Unix 纳秒）
}

// NewLimiter 创建限流器，初始速率为 0。
func NewLimiter() *Limiter {
	return &Limiter{
		last: time.Now().UnixNano(),
	}
}

// SetRate 动态调整目标速率，变更立即生效。
func (l *Limiter) SetRate(qps int) {
	atomic.StoreInt64(&l.rate, int64(qps))
}

// Allow 尝试获取一个令牌。若桶内无可用令牌则立即返回 false，不阻塞。
func (l *Limiter) Allow() bool {
	rate := atomic.LoadInt64(&l.rate)
	if rate == 0 {
		return false
	}

	now := time.Now().UnixNano()
	last := atomic.LoadInt64(&l.last)
	elapsed := now - last

	if elapsed > 0 {
		increment := elapsed * rate / 1e9
		if increment > rate {
			increment = rate
		}

		// 赢得 CAS 的协程负责补充令牌，其余协程直接尝试消费
		if atomic.CompareAndSwapInt64(&l.last, last, now) {
			for {
				curr := atomic.LoadInt64(&l.tokens)
				next := curr + increment
				if next > rate {
					next = rate
				}
				if atomic.CompareAndSwapInt64(&l.tokens, curr, next) {
					break
				}
			}
		}
	}

	for {
		curr := atomic.LoadInt64(&l.tokens)
		if curr <= 0 {
			return false
		}
		if atomic.CompareAndSwapInt64(&l.tokens, curr, curr-1) {
			return true
		}
	}
}
