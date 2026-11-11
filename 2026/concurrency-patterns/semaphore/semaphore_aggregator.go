// 🔔 场景背景
//
// 在电商或金融系统的“商品详情页”或“用户资产页”接口中，通常需要调用多个下游微服务（如：商品基础信息、库存服务、促销优惠、评价服务、推荐服务等）来拼装最终的数据返回给前端。
// 为了降低接口响应时间（RT），我们不能串行调用，而需要并发调用。但如果不加限制，几万 QPS 冲进来会瞬间把下游微服务压垮。
// 这个综合代码结合了以下生产级要素：
//
// 1、信号量限流：限制同时请求下游服务的最大并发 Goroutine 数。
// 2、Context 超时与取消机制：如果整个聚合过程超过 800ms，或者用户取消请求，立即中断所有未完成的下游调用，防止资源浪费。
// 3、并发安全与结果收集：使用 sync.Mutex 或并发安全的数据结构收集不同下游的结果。
// 4、异常防御（Panic Recovery）：单下游服务崩溃不影响其他下游数据组装。
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProductDetail 聚合的商品详情页面数据结构
type ProductDetail struct {
	ProductID string            `json:"product_id"`
	BaseInfo  string            `json:"base_info"` // 基础信息
	Stock     int               `json:"stock"`     // 库存
	Promotion string            `json:"promotion"` // 优惠活动
	Reviews   []string          `json:"reviews"`   // 热门评价
	Errors    map[string]string `json:"errors"`    // 记录哪些下游服务降级/超时
	mu        sync.Mutex        `json:"-"`         // 保护并发写入 ProductDetail 的互斥锁
}

// DownstreamTask 下游调用的任务结构
type DownstreamTask struct {
	Name string
	Exec func(ctx context.Context) (interface{}, error)
}

// 模拟 RPC 远程调用
func mockRPCCall(ctx context.Context, duration time.Duration, mockResult interface{}) (interface{}, error) {
	select {
	case <-time.After(duration): // 模拟网络延迟正常返回
		return mockResult, nil
	case <-ctx.Done(): // 模拟调用过程中，上层 context 超时取消
		return nil, ctx.Err()
	}
}

// 模拟抛出异常的服务
func mockPanicCall(ctx context.Context) (interface{}, error) {
	time.Sleep(50 * time.Millisecond)
	panic("数据库连接池空指针 (NIL Pointer Dereference)")
}

// 并发安全的错误记录方法
func (d *ProductDetail) recordError(serviceName string, errMsg string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Errors[serviceName] = errMsg
}

// 并发安全的数据填充方法
func (d *ProductDetail) fillData(serviceName string, data interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch serviceName {
	case "BaseInfoService":
		d.BaseInfo = data.(string)
	case "StockService":
		d.Stock = data.(int)
	case "PromotionService":
		d.Promotion = data.(string)
	}
}

func main() {
	// 模拟用户请求详情页面，主超时时间设置为 800ms
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	productID := "P_9527"
	detail := ProductDetail{
		ProductID: productID,
		Errors:    make(map[string]string),
	}

	// 定义需要并发调用的下游微服务列表（共 5 个服务）
	tasks := []DownstreamTask{
		{
			Name: "BaseInfoService",
			Exec: func(ctx context.Context) (interface{}, error) {
				// 模拟耗时 200ms
				return mockRPCCall(ctx, 200*time.Millisecond, "iPhone 20 Pro Max 256G")
			},
		},
		{
			Name: "StockService",
			Exec: func(ctx context.Context) (interface{}, error) {
				// 模拟耗时 150ms
				return mockRPCCall(ctx, 150*time.Millisecond, 100)
			},
		},
		{
			Name: "PromotionService",
			Exec: func(ctx context.Context) (interface{}, error) {
				// 模拟超时场景：耗时 1200ms（超过了 Context 的 800ms）
				return mockRPCCall(ctx, 1200*time.Millisecond, "满5000减600券")
			},
		},
		{
			Name: "ReviewService",
			Exec: func(ctx context.Context) (interface{}, error) {
				// 模拟异常触发场景（Panic）
				return mockPanicCall(ctx)
			},
		},
		{
			Name: "RecommendationService",
			Exec: func(ctx context.Context) (interface{}, error) {
				// 模拟耗时 300ms
				return mockRPCCall(ctx, 300*time.Millisecond, []string{"无线充电器", "手机壳"})
			},
		},
	}

	// 配置信号量限流器：下游最大并发数限制为 2 （模拟下游保护机制）
	maxDownstreamConcurrency := 2
	sem := make(chan struct{}, maxDownstreamConcurrency)

	var wg sync.WaitGroup

	startTime := time.Now()
	fmt.Printf("🚀 [Aggregator] 开始并发聚合商品 %s 的数据 (最大并发限制: %d)...\n", productID, maxDownstreamConcurrency)

	// 循环分发并发任务
	for _, task := range tasks {
		wg.Add(1)

		// 核心：在主线程获取信号量，拿不到令牌则阻塞，严格防止创建超额协程
		select {
		case sem <- struct{}{}: // 成功拿到信号量令牌
		case <-ctx.Done(): // 主 Context 已经超时，直接跳过后续任务提交
			fmt.Printf("⚠️ [Aggregator] 主流程已超时/取消，取消提交任务: %s\n", task.Name)
			detail.recordError(task.Name, ctx.Err().Error())
			wg.Done()
			continue
		}

		// 启动并发协程执行任务
		go func(t DownstreamTask) {
			defer wg.Done()

			// 任务完成后，释放信号量归还令牌
			defer func() { <-sem }()

			// 增强防御性编程：防止单个任务 Panic 击垮聚合器主协程
			defer func() {
				if r := recover(); r != nil {
					errMsg := fmt.Sprintf("panic recovered: %v", r)
					fmt.Printf("❌ [%s] 执行过程中崩溃: %s\n", t.Name, errMsg)
					detail.recordError(t.Name, errMsg)
				}
			}()

			// 执行下游调用，传入主 Context 以便感知超时/取消
			res, err := t.Exec(ctx)
			if err != nil {
				fmt.Printf("❌ [%s] 调用失败: %v\n", t.Name, err)
				detail.recordError(t.Name, err.Error())
				return
			}

			// 组装成功结果
			fmt.Printf("✅ [%s] 调用成功!\n", t.Name)
			detail.fillData(t.Name, res)
		}(task)
	}

	wg.Wait()

	fmt.Printf("\n⏱️ [Aggregator] 聚合完成！总耗时: %v\n", time.Since(startTime))
	fmt.Println("--------------------------------------------------")
	fmt.Printf("📦 基础数据: %s\n", detail.BaseInfo)
	fmt.Printf("📦 实时库存: %d\n", detail.Stock)
	fmt.Printf("📦 优惠信息: %s (降级处理)\n", detail.Promotion)
	fmt.Printf("❌ 失败/超时服务汇总: %+v\n", detail.Errors)
	fmt.Println("--------------------------------------------------")
}
