package validator

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// InputType 定义输入源类型常量
type InputType int

const (
	TypeUnknown InputType = iota
	TypeLocalFile
	TypeRemoteFile
)

// ValidateInput 校验输入参数并返回其类型
// 若输入为空或格式无效，返回对应的错误信息
func ValidateInput(input string) (InputType, error) {
	if input == "" {
		return TypeUnknown, fmt.Errorf("input cannot be empty")
	}

	input = strings.TrimSpace(input)

	// 判断是否为网络 URL
	if isNetworkURL(input) {
		return TypeRemoteFile, nil
	}

	// 判断是否为本地文件
	return validateLocalFile(input)
}

// isNetworkURL 检查字符串是否以 http:// 或 https:// 开头
// 采用前缀匹配而非完整 URL 解析，保持简单高效
func isNetworkURL(input string) bool {
	lower := strings.ToLower(input)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

// validateLocalFile 检查本地文件是否存在且可读
// 返回 TypeLocalFile 或错误信息
func validateLocalFile(path string) (InputType, error) {
	// 检查文件是否存在
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TypeUnknown, fmt.Errorf("file not found: %s", path)
		}
		return TypeUnknown, fmt.Errorf("failed to access file: %w", err)
	}

	// 检查是否为目录而非文件
	if info.IsDir() {
		return TypeUnknown, fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// 检查文件是否可读
	file, err := os.Open(path)
	if err != nil {
		return TypeUnknown, fmt.Errorf("file is not readable: %w", err)
	}
	err = file.Close()
	if err != nil {
		return 0, err
	}

	return TypeLocalFile, nil
}

// IsValidURL 对 URL 进行更严格的格式校验（可选，供外部调用）
// 使用 net/url 解析器验证 URL 结构的合法性
func IsValidURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
