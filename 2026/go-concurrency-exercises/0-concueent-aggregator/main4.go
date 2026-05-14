//////////////////////////////
// AI 创作系统的“排队机”
//
// 背景:
// Azure 账号有频率限制（比如每分钟只能请求 3 次）
// 如果你开了 10 个协程同时去调，瞬间就会报 429 Too Many Requests。
// 你需要一个“漏斗”，不管上层有多少协程在催，下层必须按照固定的节奏（比如每 2 秒一个）去调用 API。
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 方案一：用 time.Tick 实现“匀速执行”
	//ticker := time.NewTicker(2 * time.Second)
	//defer ticker.Stop()
	//
	//tasks := []string{"任务1", "任务2", "任务3", "任务4"}
	//
	//for _, task := range tasks {
	//	fmt.Printf("[%s] 正在等待获取 API 调用令牌...\n", time.Now().Format("15:04:05"))
	//
	//	<-ticker.C // 👈 关键：这里会阻塞，直到 2 秒钟的信号传过来
	//
	//	fmt.Printf("[%s] 🚀 令牌拿到！开始生成: %s\n", time.Now().Format("15:04:05"), task)
	//	time.Sleep(500 * time.Millisecond)
	//}

	// 方案二：带“缓冲”的并发控制（漏桶算法）
	// 定义并发池，限制只能有 3 个任务同时在跑
	limit := make(chan struct{}, 3)

	var wg sync.WaitGroup
	tasks := 10

	for i := 0; i < tasks; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 尝试站位：如果 limit 满了，阻塞
			limit <- struct{}{}

			fmt.Printf("[%s] 📖 章节 %d 开始调用 AI 生成...\n", time.Now().Format("15:04:05"), id)
			time.Sleep(2 * time.Second)
			fmt.Printf("[%s] ✅ 章节 %d 生成完毕\n", time.Now().Format("15:04:05"), id)

			// 释放位置：任务结束，腾出空间给后面任务
			<-limit
		}(i)
	}

	wg.Wait()
	fmt.Println("OK")
}
