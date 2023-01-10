package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	secret = "pwQ7nVkAQ7d9oTXxYZ4syTcvueaCDIZr44XJCf70H3pLkL02C3pC3cacmNkVygo2"
)

// 校验签名
func Verify(signParams map[string]interface{}, sign string) bool {
	if _, ok := signParams["signature"]; ok {
		delete(signParams, "signature")
	}

	// 按照ASCII码从小到大排序，将键名存储在切片中
	var fields []string
	for key := range signParams {
		fields = append(fields, key)
	}
	sort.Strings(fields)

	// 拼接字符串，如：field1=value1&field2=value2&...
	var buf bytes.Buffer
	for key, field := range fields {
		if field != "" && signParams[field] != "" {
			val := ""
			switch signParams[field].(type) {
			case int:
				val = strconv.Itoa(signParams[field].(int))
			case string:
				val = signParams[field].(string)
			}

			if key != len(fields)-1 {
				val += "&"
			}
			buf.WriteString(field + "=" + val)
		}
	}

	// sha256加密
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(buf.String()))
	signResult := hex.EncodeToString(hash.Sum(nil))

	return sign == strings.ToUpper(signResult)
}

func main() {
	params := map[string]interface{}{
		"app-version":  "1.0",
		"app-key":      "21e10936477411ec881a40b076627a40",
		"timestamp":    1638441452,
		"nonce":        "LO5Y5aSRczn16xFR6z",
		"client":       "iOS",
		"uuid":         "93E1FDB5-BD6D-4CE3-AC2C-A901ACBA7938",
		"lang":         "EN",
		"access-token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJiY19hcGkuaHJhY3R1YWwuY29tIiwiYXVkIjoiOTNFMUZEQjUtQkQ2RC00Q0UzLUFDMkMtQTkwMUFDQkE3OTM4IiwiaWF0IjoxNjM4NDMzMzU4LCJuYmYiOjE2Mzg0MzMzNjAsImV4cCI6MTYzODQ0MDU1OCwiZGF0YSI6eyJ1dWlkIjoiOTNFMUZEQjUtQkQ2RC00Q0UzLUFDMkMtQTkwMUFDQkE3OTM4IiwiYXBwX2tleSI6IjIxZTEwOTM2NDc3NDExZWM4ODFhNDBiMDc2NjI3YTQwIiwiY2xpZW50IjoiaU9TIiwibGFuZyI6IkVOIiwiaXAiOiIxMDMuOTcuMjAxLjQifX0.l1HdTqKGiBYXLPgnEUKtCq8ztimOlezUPx9RPMKCVPM",
		"api_url":      "f1335bdd92f72e8d8366426d5d697869",
		"signature":    "",
	}
	sign := "534EE7575020CAF7E62F2288F2E9A67BC7730A0FCC4A0DB7B1903239A0DD009A"
	result := Verify(params, sign)
	fmt.Println(result)
}
