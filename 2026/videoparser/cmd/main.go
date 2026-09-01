package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"toolbox/2026/videoparser/internal/parser"
)

func main() {
	var input string
	var showVersion bool

	flag.StringVar(&input, "i", "", "输入视频文件路径或网络 URL（http/https）")
	flag.BoolVar(&showVersion, "version", false, "显示 ffprobe 版本信息")
	flag.Usage = printUsage
	flag.Parse()

	if showVersion {
		printVersion()
		return
	}

	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	result := parser.Parse(input)

	outputJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to marshal result: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(outputJSON))

	if !result.Success {
		os.Exit(1)
	}
}

func printVersion() {
	version, err := parser.GetFFprobeVersion()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get ffprobe version: %v\n", err)
		os.Exit(1)
	}

	output := map[string]string{
		"ffprobe_version": version,
		"status":          "available",
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

func printUsage() {
	fmt.Println("Video Parser CLI Tool v1.0")
	fmt.Println("")
	fmt.Println("高性能视频信息解析工具，基于 FFmpeg 实现")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  videoparser -i <输入>    解析视频信息")
	fmt.Println("  videoparser --version   显示 ffprobe 版本")
	fmt.Println("")
	fmt.Println("输入类型:")
	fmt.Println("  本地文件: videoparser -i video.mp4")
	fmt.Println("  网络链接: videoparser -i https://example.com/video.mp4")
	fmt.Println("")
	fmt.Println("选项:")
	flag.PrintDefaults()
}
