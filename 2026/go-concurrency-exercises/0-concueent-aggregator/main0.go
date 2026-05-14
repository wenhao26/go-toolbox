//////////////////////////////
// 并行的聚合查询
//
// 常规做法（串行）：耗时 = 查询用户 ($100ms$) + 查询订单 ($200ms$) + 查询余额 ($150ms$) = $450ms
// 并发做法：耗时 = $max(100ms, 200ms, 150ms)$ = $200ms$
// 结论：总耗时由最慢的那个接口决定，性能提升显著
//

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// UserData 用户聚合数据
type UserData struct {
	Info    string
	Orders  []string
	Balance float64
}

// GetUserAggregatedData 获取用户聚合数据
func GetUserAggregatedData(ctx context.Context, userID int64) (*UserData, error) {
	var wg sync.WaitGroup
	data := &UserData{}

	// 用于捕获协程内部发生的错误
	errChan := make(chan error, 3)

	// 开启协程查询用户信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TODO 模拟数据查询
		time.Sleep(100 * time.Millisecond)

		//info, err := db.QueryInfo(userID)
		//if err != nil {
		//	errChan <- fmt.Errorf("查询信息失败: %v", err) // 具体的错误收集在这里！
		//	return
		//}

		data.Info = fmt.Sprintf("艾塔娜-Aitana-%d", userID)
		log.Println("Info 查询完成")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(4 * time.Second)
		data.Orders = []string{"D123456", "E123456"}
		log.Println("Orders 查询完成")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		data.Balance = 100.57
		log.Println("Balance 查询完成")
	}()

	// 用来通知主流程是否全部完成
	allDone := make(chan struct{})

	// 等待所有查询完成
	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
		return data, nil // 先收到信号（被关闭）
	case <-ctx.Done():
		log.Println("执行超时")
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
		//default:
		//	return data, nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userID := int64(9527)
	result, err := GetUserAggregatedData(ctx, userID)
	if err != nil {
		fmt.Printf("聚合查询失败: %v\n", err)
		return
	}

	fmt.Printf("最终结果: %+v\n", result)
}
