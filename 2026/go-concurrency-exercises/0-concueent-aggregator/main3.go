//////////////////////////////
// 高并发下的“英雄救美” —— SingleFlight 模式
//
// 背景:
// 假设你的“统计松鼠”或“内容页”突然爆火。1000 个用户在同一毫秒请求同一个 ID 的文章。
// - 初级思维：1000 个协程全部去查数据库。结果：数据库瞬间宕机。
// - 中级思维：加缓存。但如果缓存刚好过期，这 1000 个请求还是会瞬间穿透到数据库。
// - ✅架构师思维（Go 的魅力）：既然是查同一个 ID，为什么要跑 1000 次？
//     我让第一个人去查，剩下的 999 个人在门口等着，等第一个人带回结果，大家共享。
//
// 这种模式在 Go 中被称为 SingleFlight（单飞变合飞）。
// 手写一个 SingleFlight 调度器
// 我们不直接用官方库，先用原生代码实现核心逻辑，让你感受 Go 处理这种问题的优雅
//

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Result 包装返回结果
type Result struct {
	Data string
	Err  error
}

// Call 代表一个正在飞行（进行中）的请求
type Call struct {
	wg  sync.WaitGroup
	res Result
}

// Group 调度中心
type Group struct {
	mu sync.Mutex       // 并发锁（保护 m）
	m  map[string]*Call // 存放正在执行的请求
}

// Do 核心方法：相同的 key，无论并发多少次，实际逻辑只执行一次
func (g *Group) Do(key string, fn func(ctx context.Context) (string, error)) (string, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}

	// 如果发现有人已经在查这个 key
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 等待，直接当前请求返回
		return c.res.Data, c.res.Err
	}

	// 如果是首个请求
	c := new(Call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 去查数据
	// 无论外部传进来的 ctx 是什么，首个用户必须保证自己不卡死
	c.res.Data, c.res.Err = func() (string, error) {
		// 强制给底层 IO 设一个保底超时（例如 5 秒）
		innerCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return fn(innerCtx)
	}()

	// 广播给所有请求（告诉等待中请求，已经完成了）
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) // 完成后，从“飞行列表”删除
	g.mu.Unlock()

	return c.res.Data, c.res.Err
}

func main() {
	g := &Group{}

	const key = "key_9527"

	var wg sync.WaitGroup

	// TODO 模拟 10 个并发同时请求
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Printf("用户 %d 发起请求...\n", id)

			// 使用 SingleFlight 执行
			data, _ := g.Do(key, func(ctx context.Context) (string, error) {
				log.Printf("👉👉👉！！！用户 %d 正在执行数据库查询（这行应该只出现一次）！！！\n", id)
				// time.Sleep(3 * time.Second)

				select {
				case <-time.After(4 * time.Second):
					return "响应内容", nil
				case <-ctx.Done():
					return "", ctx.Err() // 首个用户请求超时了
				}
			})

			log.Printf("用户 %d 拿到结果: %s\n", id, data)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有并发请求已处理完毕。")
}
