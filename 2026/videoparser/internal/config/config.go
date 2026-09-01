package config

import (
	"os"
	"strconv"
	"time"
)

var (
	// FFprobeTimeout 定义 ffprobe 命令执行的超时时间，防止卡死
	// 可通过环境变量 VIDEO_PARSER_TIMEOUT 覆盖，单位：秒
	FFprobeTimeout = getEnvDuration("VIDEO_PARSER_TIMEOUT", 60)

	// FFprobeAnalyzeDuration 控制 ffprobe 分析文件的时间长度（微秒）
	// 较小的值可提升启动速度，但对某些格式可能遗漏元数据。
	// 可通过环境变量 VIDEO_PARSER_ANALYZE_DURATION 覆盖，单位：微秒
	FFprobeAnalyzeDuration = getEnvString("VIDEO_PARSER_ANALYZE_DURATION", "5000000")

	// FFprobeProbeSize 控制 ffprobe 读取的数据量（字节），用于元数据分析
	// 可通过环境变量 VIDEO_PARSER_PROBE_SIZE 覆盖，单位：字节
	FFprobeProbeSize = getEnvString("VIDEO_PARSER_PROBE_SIZE", "5000000")
)

// getEnvDuration 从环境变量读取超时时间（秒），若失败或未设置则返回默认值
func getEnvDuration(key string, defaultValue int) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if seconds, err := strconv.Atoi(val); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return time.Duration(defaultValue) * time.Second
}

// getEnvString 从环境变量读取字符串配置，若未设置则返回默认值
func getEnvString(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}
