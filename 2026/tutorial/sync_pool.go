package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

var newCount int32 // 统计到底调用了多少次 New

// bufferPool 定义一个全局的 Buffer 池
var bufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("【GC/Pool提示】正在分配一个新的 bytes.Buffer 对象...")
		atomic.AddInt32(&newCount, 1)
		return new(bytes.Buffer)
	},
}

// - 实现高性能对象复用（减少 GC 压力）
// 场景：在高并发后端服务中（比如处理 JSON 解析、Buffer 拼接），如果每个请求都频繁地 make([]byte, 1024) 或创建临时大对象，
// 会导致内存剧烈波动，给 Go 的垃圾回收（GC）带来巨大压力。
//
// 惊艳之处：sync.Pool 是 Go 标准库自带的临时对象缓存池。它可以让对象被循环利用，显著降低内存分配开销和 GC 停顿（Stop-The-World）。
//
// 学习收获：理解内存分配的底层瓶颈。在写高吞吐（High-Throughput）中间件或微服务时，善用 sync.Pool 是从“普通开发者”迈向“高级工程优化”的重要一步。

func main() {
	var wg sync.WaitGroup

	workerNum := 10
	loops := 100 // 每个协程循环 100 次

	// 模拟 10 个并发协程同时高频使用 Pool
	for i := 0; i < workerNum; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < loops; j++ {
				// 1、获取
				buf := bufferPool.Get().(*bytes.Buffer)

				// 2、使用
				buf.WriteString("some data")

				// 3、重置并放回
				buf.Reset()
				bufferPool.Put(buf)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("总共执行了 %d 次操作\n", workerNum*loops)
	fmt.Printf("实际触发 New 创建新对象的次数: %d\n", newCount)
}
