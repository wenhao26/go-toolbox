// 并发处理：用户行为埋点，统计抽奖人数，抽奖记录写入，抽奖成功写入奖励记录
// 任务完成才返回抽奖结果

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// 定义任务接口
type Task interface {
	Do() error
}

// UserBehavior 用户行为埋点
type UserBehavior struct {
	UserID int
	Action string
}

func (t *UserBehavior) Do() error {
	time.Sleep(4 * time.Second)
	fmt.Printf(" - Collect user behavior:user-id=%d,action=%s...\n", t.UserID, t.Action)
	return nil
}

// UserLotteryCount 统计抽奖人数
type UserLotteryCount struct{}

func (t *UserLotteryCount) Do() error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println(" - Statistical lottery...")
	return nil
}

// UserLotteryRecord 抽奖记录
type UserLotteryRecord struct {
	UserID    int
	PrizeID   int
	CreatedAt string
}

func (t *UserLotteryRecord) Do() error {
	time.Sleep(3 * time.Second)
	fmt.Printf(" - Write lottery record:user-id=%d,prize-id=%d,created_at=%s...\n", t.UserID, t.PrizeID, t.CreatedAt)
	return nil
}

// UserRewardRecord 奖励记录
type UserRewardRecord struct {
	UserID    int
	PrizeID   int
	RewardID  int
	Award     string
	CreatedAt string
}

func (t *UserRewardRecord) Do() error {
	time.Sleep(3 * time.Second)
	fmt.Printf(
		" - Write reward record:user-id=%d,prize-id=%d,reward-id=%d,award=%s,created_at=%s...\n",
		t.UserID,
		t.PrizeID,
		t.RewardID,
		t.Award,
		t.CreatedAt,
	)
	return nil
}

// 抽奖
func DrawLottery(userID, prizeID, rewardID int, award string) error {
	var wg sync.WaitGroup

	// 创建任务列表
	tasks := []Task{
		&UserBehavior{
			UserID: userID,
			Action: award,
		},
		&UserLotteryCount{},
		&UserLotteryRecord{
			UserID:    userID,
			PrizeID:   prizeID,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		&UserRewardRecord{
			UserID:    userID,
			PrizeID:   prizeID,
			RewardID:  rewardID,
			Award:     award,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}

	// 创建线程池
	pool, err := ants.NewPool(4)
	if err != nil {
		return fmt.Errorf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	// 提交任务到线程池
	for _, task := range tasks {
		wg.Add(1)
		taskCopy := task // 避免循环变量问题
		err := pool.Submit(func() {
			defer wg.Done()
			if err := taskCopy.Do(); err != nil {
				fmt.Printf("任务执行失败: %v\n", err)
			}
		})
		if err != nil {
			return fmt.Errorf("提交任务失败: %v", err)
		}
	}

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("所有任务已完成，返回抽奖结果...")
	return nil
}

func main() {
	fmt.Println("开始用户抽奖流程...")
	userID := 1688
	prizeID := 202501170100001
	rewardID := 202501170200001
	award := "抽中100元优惠券"

	startTime := time.Now()
	if err := DrawLottery(userID, prizeID, rewardID, award); err != nil {
		fmt.Printf("抽奖流程失败: %v\n", err)
	} else {
		fmt.Printf("抽奖流程成功完成，耗时: %v\n", time.Since(startTime))
	}
}
