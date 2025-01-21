// 使用ants协程池来并发处理任务
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// User 用户信息结构体
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ProcessUserData 处理用户数据
func ProcessUserData(userID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取用户数据
	user, err := FetchUserData(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 模拟数据处理
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("处理用户数据: ID=%d, Name=%s, Username=%s, Email=%s\n", user.ID, user.Name, user.Username, user.Email)
}

// 获取用户数据
func FetchUserData(userID int) (*User, error) {
	// 模拟 HTTP 请求
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求用户 %d 数据失败: %v", userID, err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取用户 %d 数据失败: %v", userID, err)
	}

	// 解析 JSON 数据
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("解析用户 %d 数据失败: %v", userID, err)
	}

	return &user, nil
}

func main() {
	var wg sync.WaitGroup

	// 定义用户ID列表
	userIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 创建一个协程池，大小为4
	pool, err := ants.NewPool(6)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	startTime := time.Now()

	// 提交任务到协程池
	for _, userID := range userIDs {
		wg.Add(1)
		userIDCopy := userID // 避免闭包捕获循环变量
		_ = pool.Submit(func() {
			ProcessUserData(userIDCopy, &wg)
		})
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("所有任务已完成，总耗时: %v\n", elapsedTime)
}
