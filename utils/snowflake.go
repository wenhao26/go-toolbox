package utils

import (
	"sync"
	"time"
)

type Worker struct {
	// 互斥锁，确保并发安全
	mux sync.Mutex
	// 记录时间戳
	timestamp int64
	// 当前毫秒已经生成的id序列号（从0开始累加），1毫秒内最多生成的4096个ID
	number int64
}

func NewWorker() *Worker {
	return &Worker{timestamp: 0, number: 0}
}

func (w *Worker) GetID() int64 {
	// 设置去年今天的时间戳
	epoch := int64(1645611194)
	idLength := uint(9)

	w.mux.Lock()
	defer w.mux.Unlock()

	now := time.Now().Unix()
	if w.timestamp == now {
		w.number++
		// 此处为最大节点ID，大概是2^9-1条
		if w.number > (-1 ^ (-1 << idLength)) {
			for now <= w.timestamp {
				now = time.Now().Unix()
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	return (now-epoch)<<idLength | (int64(1) << 1) | (w.number)
}
