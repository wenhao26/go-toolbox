package main

import (
	"fmt"
	"log"
	"time"

	"toolbox/2025/httpclient/client"
)

// PostsResponse 定义 posts 响应结构
type PostsResponse struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	// 创建 HTTP 客户端，设置10秒超时
	c := client.NewHTTPClient(10 * time.Second)

	// 定义请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 示例一：发送 GET 请求
	var getPostsResponse PostsResponse
	err := c.Get("https://jsonplaceholder.typicode.com/posts/1", headers, &getPostsResponse)
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	// fmt.Println(getPostsResponse.UserID, getPostsResponse.ID, getPostsResponse.Title, getPostsResponse.Body)
	fmt.Println("GET Response:", getPostsResponse)

	// 示例二：发送 POST 请求
	postBody := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	var postPostsResponse PostsResponse
	err = c.Post("https://jsonplaceholder.typicode.com/posts", headers, postBody, &postPostsResponse)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Response:", postPostsResponse)

	// 示例 3: 发送 PUT 请求
	putBody := map[string]interface{}{
		"id":     1,
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	var putPostsResponse PostsResponse
	err = c.Put("https://jsonplaceholder.typicode.com/posts/1", headers, putBody, &putPostsResponse)
	if err != nil {
		log.Fatalf("PUT request failed: %v", err)
	}
	fmt.Println("PUT Response:", putPostsResponse)

	// 示例 4: 发送 DELETE 请求
	var deletePostsResponse PostsResponse
	err = c.Delete("https://jsonplaceholder.typicode.com/posts/1", headers, &deletePostsResponse)
	if err != nil {
		log.Fatalf("DELETE request failed: %v", err)
	}
	fmt.Println("DELETE Response:", deletePostsResponse)

	// 示例 5: 重试机制的 GET 请求
	retries := 3
	var retryResponse PostsResponse
	_, err = c.DoRequestWithRetry("GET", "https://jsonplaceholder.typicode.com/posts/1", headers, nil, retries)
	if err != nil {
		log.Fatalf("GET with retry failed: %v", err)
	}
	fmt.Println("GET with retry Response:", retryResponse)

}
