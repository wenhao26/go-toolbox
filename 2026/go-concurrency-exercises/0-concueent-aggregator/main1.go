//////////////////////////////
// 电商系统的“用户风控画像”
//
// 业务需求：
// - 第一步（串行）：必须先根据手机号去数据库查询用户ID。如果查不到，后面全不用做了。
//
// - 第二步（并发）：拿到用户ID后，需要同时去三个不同的地方抓数据：
//	a.查询黑名单系统（如果命中黑名单，直接返回，不再等其他数据）。
//	b.查询历史消费金额（MySQL）。
//	c.查询信用评分（第三方风控接口）。
//
// - 约束：整个过程限时 3秒。
//

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// RiskTags 用户风控标签
type RiskTags struct {
	UserID        int64   // 用户ID
	IsBlacklisted bool    // 是否在黑名单
	TotalSpent    float64 // 历史消费金额
	CreditScore   float64 // 信用评分
}

// 查询用户
func getUserID(phoneNumber string) (int64, error) {
	// TODO 模拟查询耗时
	time.Sleep(200 * time.Microsecond)
	return 9527, nil
}

// 是否命中黑名单检查
func hitBlacklist(userID int64) bool {
	// TODO 模拟去黑名单系统查询数据
	time.Sleep(1 * time.Second)
	log.Printf("用户ID %d 黑名单结果查询完成\n", userID)
	return false
}

// 查询历史总消耗金额
func getTotalSpent(userID int64) float64 {
	// TODO 模拟去交易系统查询数据
	time.Sleep(1 * time.Second)
	log.Printf("用户ID %d 总消耗金额结果查询完成\n", userID)
	return 5001.45
}

// 查询信用评分
func getCreditScore(userID int64) float64 {
	// TODO 模拟去交易系统查询数据
	time.Sleep(400 * time.Microsecond)
	log.Printf("用户ID %d 信用评分结果查询完成\n", userID)
	return 3.5
}

func getRiskProfile(ctx context.Context, userID int64) (*RiskTags, error) {
	var wg sync.WaitGroup

	errChan := make(chan error, 3)
	riskTags := &RiskTags{
		UserID: userID,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		riskTags.IsBlacklisted = hitBlacklist(userID)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		riskTags.TotalSpent = getTotalSpent(userID)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		riskTags.CreditScore = getCreditScore(userID)
	}()

	allDone := make(chan struct{})

	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
		if riskTags.IsBlacklisted {
			return riskTags, fmt.Errorf("该用户已命中黑名单")
		}
		return riskTags, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return riskTags, err
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	phoneNumber := "+1 (646) 555-0143"

	userID, err := getUserID(phoneNumber)
	if err != nil {
		log.Printf("查询异常: %v", err)
		return
	}
	if userID == 0 {
		log.Println("未找到该用户")
		return
	}

	result, err := getRiskProfile(ctx, userID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", result)
}
