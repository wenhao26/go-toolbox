package utils

import (
	"os/exec"
	"runtime"
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