package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Notifier 通知器结构体
type Notifier struct {
	webHook   string
	secretKey string
}

// NewNotifier 创建通知器实例
func NewNotifier(webHook, secretKey string) *Notifier {
	return &Notifier{
		webHook:   webHook,
		secretKey: secretKey,
	}
}

// sign 签名
func (n *Notifier) Sign(timestamp int64) string {
	if n.secretKey == "" {
		return ""
	}

	keyString := fmt.Sprintf("%d\n%s", timestamp, n.secretKey)
	hash := hmac.New(sha256.New, []byte(keyString))
	signature := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature)
}

// Send 发送
func (n *Notifier) Send(message Message) map[string]interface{} {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 配置 Transport 以忽略 SSL 证书验证
	}
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // 增加超时配置，可选
		Transport: transport,
	}

	payloadBytes, err := json.Marshal(message.Payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", n.webHook, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {
		var responseData map[string]interface{}

		jsonErr := json.Unmarshal(responseBody, &responseData)
		if jsonErr != nil {
			return map[string]interface{}{
				"err": map[string]interface{}{
					"code": 400,
					"msg":  fmt.Sprintf("JSON Decode Error: %v. Raw body: %s", jsonErr, string(responseBody)),
				},
			}
		}

		return responseData
	}

	return map[string]interface{}{
		"err": map[string]interface{}{
			"code": resp.StatusCode,
			"msg":  string(responseBody),
		},
	}
}
