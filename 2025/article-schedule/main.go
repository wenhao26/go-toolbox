// 需求：文章定时发布模拟
// 存储：Redis
// 目标：实时监控文章定时参数，实现定时发布处理
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"toolbox/2025/article-schedule/rdb"
)

var ctx = context.Background()

// Article 文章信息
type Article struct {
	Title     string
	Content   string
	DelayTime float64 // 延迟时间，单位秒
}

// addScheduleArticle 添加定时文章
func (a *Article) addScheduleArticle(title, content string, publishTime time.Time) error {
	// 将发布时间转换为Unix时间戳
	timestamp := float64(publishTime.Unix())

	article := Article{
		Title:     title,
		Content:   content,
		DelayTime: timestamp,
	}

	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("JSON 编码失败: %v", err)
	}

	// 添加数据到Redis有序集合
	redisStorage := rdb.NewRedisClient("localhost:6379", 9)
	err = redisStorage.Rdb.ZAdd(ctx, "schedule_article", &redis.Z{
		Score:  timestamp,
		Member: jsonData,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to schedule post: %v", err)
	}

	return nil
}

// publishArticle 发布文章
func publishArticle(data string) {
	// 模拟发布逻辑
	fmt.Printf("Publishing article: %s\n", data)

	// 从Redis中移除已发布的文章
	redisStorage := rdb.NewRedisClient("localhost:6379", 9)
	redisStorage.Rdb.ZRem(ctx, "schedule_article", data)
}

// monitorScheduleArticle 监控定时文档
func (a *Article) monitorScheduleArticle() {
	ticker := time.NewTicker(1 * time.Second) // 每秒检查一次
	defer ticker.Stop()

	redisStorage := rdb.NewRedisClient("localhost:6379", 9)

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()

			// 获取当前时间是否存在需要发布的文章
			articles, err := redisStorage.Rdb.ZRangeByScore(ctx, "schedule_article", &redis.ZRangeBy{
				Min: "0",             // 最小时间戳
				Max: fmt.Sprint(now), // 当前时间戳
			}).Result()
			if err != nil {
				log.Printf("Failed to fetch scheduled article: %v", err)
				continue
			}

			// 发布文章
			for _, article := range articles {
				publishArticle(article)
			}
		}
	}
}

func main() {
	article := Article{}

	// 启动监控 Goroutine
	go article.monitorScheduleArticle()

	title := "基于Go整合Gossip+WebSocket的轻量级分布式消息队列"
	content := "实现功能：订阅/发布推送、点对点推送、延时消息推送、消费失败自动重试、数据持久化、WAL预写日志、主从节点部署、集群负载均衡、ACK确认机制、内置网页端控制台。"

	// 添加一些测试文章
	_ = article.addScheduleArticle(title, content, time.Now().Add(5*time.Second))  // 5 秒后发布
	_ = article.addScheduleArticle(title, content, time.Now().Add(15*time.Second)) // 15 秒后发布
	_ = article.addScheduleArticle(title, content, time.Now().Add(25*time.Second)) // 25 秒后发布
	_ = article.addScheduleArticle(title, content, time.Now().Add(35*time.Second)) // 35 秒后发布
	_ = article.addScheduleArticle(title, content, time.Now().Add(45*time.Second)) // 45 秒后发布
	_ = article.addScheduleArticle(title, content, time.Now().Add(55*time.Second)) // 55 秒后发布

	fmt.Println("定时文章发布成功")

	// 保持主程序运行
	select {}
}
