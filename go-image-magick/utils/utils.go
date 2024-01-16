package utils

import (
	"os/exec"
	"runtime"
	"strconv"
)

// 操作系统名称
func SystemName() string {
	return runtime.GOOS
}

// 执行 ImageMagick 命令
func Magick(params []string) (string, error) {
	binPath, err := exec.LookPath("magick")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(binPath, params...)
	result, err := cmd.CombinedOutput()
	return string(result), err
}

// 去除后面精确值全是0的数值
func RemoveTrailingZeros(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 32) // 将float64转换为字符串形式，保留两位有效数字
}