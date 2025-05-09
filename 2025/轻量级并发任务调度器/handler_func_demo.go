package main

import (
	"errors"
	"fmt"
	"time"
)

//
// ==== 自定义 handler 类型 ====
//

// HandlerFunc 定义一个 函数类型
type HandlerFunc func() error

// 执行一个 handler 并捕获错误
func Execute(h HandlerFunc) {
	if err := h(); err != nil {
		fmt.Println("❌ 执行失败:", err)
	} else {
		fmt.Println("✅ 执行成功")
	}
}

//
// ==== 中间封装 ====
//

// WithLogging 一个简单的中间件：记录执行时间和日志
func WithLogging(h HandlerFunc) HandlerFunc {
	return func() error {
		fmt.Println("开始执行任务...")
		start := time.Now()

		err := h()

		fmt.Printf("🟢 执行耗时: %s\n", time.Since(start))
		fmt.Println("🔵 执行完成")
		return err
	}
}

// 另一个中间件：重试机制
func WithRetry(h HandlerFunc, retry int) HandlerFunc {
	return func() error {
		var err error
		for i := 0; i < retry; i++ {
			err = h()
			if err == nil {
				return nil
			}
			fmt.Printf("🔁 第 %d 次重试失败: %s\n", i+1, err)
		}
		return err
	}
}

//
// ==== 任务函数 ====
//

func TaskJob() error {
	fmt.Println("⚙️正在执行任务...")
	// 模拟失败
	if time.Now().Unix()%2 == 0 {
		return errors.New("模拟错误：执行失败")
	}
	return nil
}

func main() {
	// 原始 handler
	baseHandler := TaskJob

	// 日志中间件
	logged := WithLogging(baseHandler)

	// 重试中间件（最多重试2次）
	retryHandler := WithRetry(logged, 2)

	// 执行
	Execute(retryHandler)
}
