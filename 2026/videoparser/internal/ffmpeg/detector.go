package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// FFprobePath 存储检测到的 ffprobe 可执行文件绝对路径
var FFprobePath string

// init 在包被导入时自动检测 ffprobe 是否可用
func init() {
	FFprobePath = detectFFprobe()
}

// detectFFprobe 在系统 PATH 中搜索 ffprobe 可执行文件
// 根据操作系统自动添加后缀（Windows 下为 .exe）
// 若未找到则返回空字符串
func detectFFprobe() string {
	var cmdName string

	if runtime.GOOS == "windows" {
		cmdName = "ffprobe.exe"
	} else {
		cmdName = "ffprobe"
	}

	// 检查 PATH 中的 ffprobe
	if path, err := exec.LookPath(cmdName); err == nil {
		return path
	}

	// 若 PATH 查找失败，尝试常见安装路径（仅限 Windows）
	if runtime.GOOS == "windows" {
		commonPaths := []string{
			`D:\Program Files\ffmpeg\bin\ffprobe.exe`,
		}
		for _, p := range commonPaths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

// IsAvailable 检查 ffprobe 是否已成功检测到
func IsAvailable() bool {
	return FFprobePath != ""
}

// GetVersion 获取 ffprobe 版本号，用于调试或兼容性判断
// 返回字符串和可能的错误，若 ffprobe 不可用则返回空字符串
func GetVersion() (string, error) {
	if !IsAvailable() {
		return "", fmt.Errorf("ffprobe not found in system PATH")
	}

	cmd := exec.Command(FFprobePath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute ffprobe: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}
	return "", fmt.Errorf("unexpected ffprobe output")
}
