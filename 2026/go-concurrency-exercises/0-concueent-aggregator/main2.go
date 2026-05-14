//////////////////////////////
// 高可用的“内容详情页”聚合
//
// 核心与非核心业务的降级与容错：
// 在真实的微服务架构中，一个接口背后可能聚合了 10 个服务。
// 如果其中 2 个非核心服务（比如“推荐广告”或“勋章墙”）挂了或者超时了，
// 你不能让整个接口报错，而是应该“弃车保帅”——返回核心数据，漏掉非核心数据。
//
// 业务需求：
// - 核心业务（必须成功）：文章正文（MySQL）。如果这个失败，直接报错。
//
// - 次要业务（可降级）：相关推荐（Redis）、用户点赞状态（另一个微服务）。
//
// - 约束：
//	a.限时 500ms。
//  b.即使“相关推荐”或“点赞状态”查询失败（报错或超时），也要返回“文章正文”，只是对应的字段为空。
//
// 深入提升：使用 errgroup 处理部分容错
//

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ContentPage 内容页结构
type ContentPage struct {
	ID      int64    // 文章唯一标识 ID
	Article string   // 文章正文，必须返回
	Likes   int      // 点赞数
	Suggest []string // 推荐列表
}

// GetContentDetail 获取内容页详情数据
func GetContentDetail(ctx context.Context, id int64) (*ContentPage, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	contentPage := &ContentPage{
		ID: id,
	}

	// 核心任务
	coreErrChan := make(chan error, 1) // 文章正文 (只要它报错，就认为整个聚合失败)

	wg.Add(1)
	go func() {
		defer wg.Done()
		text, err := fetchArticle(ctx, id)
		if err != nil {
			coreErrChan <- err
			return
		}
		contentPage.Article = text
	}()

	// 次要任务
	wg.Add(1)
	go func() {
		defer wg.Done()
		likes, err := fetchLikes(ctx, id)
		if err != nil {
			// 【思维提升点】：非核心业务报错，记录日志，不阻塞主流程
			log.Printf("次要业务[点赞数]获取失败: %v", err)
			return
		}
		contentPage.Likes = likes
	}()

	// 次要业务
	wg.Add(1)
	go func() {
		defer wg.Done()
		suggest, err := fetchSuggest(ctx, id)
		if err != nil {
			log.Printf("次要业务[推荐列表]获取失败: %v", err)
			return
		}
		contentPage.Suggest = suggest
	}()

	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
		return contentPage, nil
	case err := <-coreErrChan:
		return nil, fmt.Errorf("获取文章内容失败: %w", err)
	case <-ctx.Done():
		// 如果超时，我们要判断核心业务是否已经拿到了
		if contentPage.Article != "" {
			log.Println("部分业务超时，但核心数据已获取，执行降级返回")
			return contentPage, nil
		}
		return nil, ctx.Err()
	}
}

// fetchArticle 获取文章正文内容
func fetchArticle(ctx context.Context, id int64) (string, error) {
	// TODO 模拟返回文章正文内容
	time.Sleep(3 * time.Second)
	content := `Early the next morning, I arrived at the hospital right on time.`
	return fmt.Sprintf("[%d]%s", id, content), nil
}

// fetchLikes 获取点赞数
func fetchLikes(ctx context.Context, id int64) (int, error) {
	// TODO 模拟 Redis 服务不可用，无法获取点赞数
	time.Sleep(100 * time.Millisecond)
	return 0, fmt.Errorf("redis service unavailable")
}

// fetchSuggest 获取推荐列表
func fetchSuggest(ctx context.Context, id int64) ([]string, error) {
	// TODO 模拟返回推荐列表
	time.Sleep(100 * time.Millisecond)
	list := []string{"文章1", "文章2", "文章3"}
	return list, nil
}

func main() {
	detail, err := GetContentDetail(context.Background(), 9527)
	if err != nil {
		log.Fatalf("接口请求失败: %v", err)
		return
	}

	fmt.Println(detail)
}
