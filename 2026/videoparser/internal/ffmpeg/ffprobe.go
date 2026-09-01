package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"toolbox/2026/videoparser/internal/config"
)

// ProbeOptions 定义 ffprobe 执行时的可选参数
type ProbeOptions struct {
	Input     string
	IsNetwork bool
}

// Probe 执行 ffprobe 命令获取视频元数据
// 返回原始 JSON 字节数据，调用方需自行解析
// 若 ffprobe 未安装或执行失败，返回错误
func Probe(opts ProbeOptions) ([]byte, error) {
	if !IsAvailable() {
		return nil, fmt.Errorf("ffprobe not found, please install FFmpeg and ensure ffprobe is in PATH")
	}

	args := buildProbeArgs(opts)

	ctx, cancel := context.WithTimeout(context.Background(), config.FFprobeTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, FFprobePath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("ffprobe execution timeout after %v", config.FFprobeTimeout)
		}
		return nil, fmt.Errorf("ffprobe execution failed: %w, stderr: %s", err, stderr.String())
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("ffprobe returned empty output, stderr: %s", stderr.String())
	}

	return stdout.Bytes(), nil
}

// buildProbeArgs 构造 ffprobe 命令行参数
// 根据是否为网络 URL 添加不同的参数组合
func buildProbeArgs(opts ProbeOptions) []string {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
	}

	if opts.IsNetwork {
		args = append(args,
			"-timeout", "30000000",
			"-user_agent", "Mozilla/5.0 (compatible; VideoParser/1.0)",
		)
	}

	args = append(args,
		"-analyzeduration", config.FFprobeAnalyzeDuration,
		"-probesize", config.FFprobeProbeSize,
		opts.Input,
	)

	return args
}
