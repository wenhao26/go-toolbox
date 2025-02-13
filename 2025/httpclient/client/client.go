// 封装 HTTP 客户端
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClient 封装的 HTTP 客户端
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建一个 HTTP 客户端
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoRequest 执行 HTTP 请求
func (c *HTTPClient) DoRequest(method, url string, headers map[string]string, body interface{}) (*http.Response, error) {
	// 序列化请求体
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 执行请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	return resp, nil
}

// ParseResponse 解析 HTTP 响应
func (c *HTTPClient) ParseResponse(resp *http.Response, result interface{}) error {
	// 确保响应体在关闭前读取
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应体
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %v", err)
		}
	}

	return nil
}

// DoRequestWithRetry 执行 HTTP 请求并支持重试机制
func (c *HTTPClient) DoRequestWithRetry(method, url string, headers map[string]string, body interface{}, retries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < retries; i++ {
		resp, err = c.DoRequest(method, url, headers, body)
		if err == nil {
			return resp, nil
		}

		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(2 * time.Second) // 重试间隔
	}

	return nil, fmt.Errorf("failed after %d retries: %v", retries, err)
}

// GET 发起 GET 请求并解析响应
func (c *HTTPClient) Get(url string, headers map[string]string, result interface{}) error {
	resp, err := c.DoRequest("GET", url, headers, nil)
	if err != nil {
		return fmt.Errorf("GET request failed: %v", err)
	}

	return c.ParseResponse(resp, result)
}

// POST 发起 POST 请求并解析响应
func (c *HTTPClient) Post(url string, headers map[string]string, body interface{}, result interface{}) error {
	resp, err := c.DoRequest("POST", url, headers, body)
	if err != nil {
		return fmt.Errorf("POST request failed: %v", err)
	}

	return c.ParseResponse(resp, result)
}

// PUT 发起 PUT 请求并解析响应
func (c *HTTPClient) Put(url string, headers map[string]string, body interface{}, result interface{}) error {
	resp, err := c.DoRequest("PUT", url, headers, body)
	if err != nil {
		return fmt.Errorf("PUT request failed: %v", err)
	}

	return c.ParseResponse(resp, result)
}

// DELETE 发起 DELETE 请求并解析响应
func (c *HTTPClient) Delete(url string, headers map[string]string, result interface{}) error {
	resp, err := c.DoRequest("DELETE", url, headers, nil)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}

	return c.ParseResponse(resp, result)
}
