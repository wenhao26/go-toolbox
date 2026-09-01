package parser

import (
	"encoding/json"
	"fmt"
	"toolbox/2026/videoparser/internal/ffmpeg"
	"toolbox/2026/videoparser/internal/model"
	"toolbox/2026/videoparser/internal/validator"
)

// ffprobeOutput 定义 ffprobe 返回的原始 JSON 结构
// 仅用于解析中间过程，不对外暴露
type ffprobeOutput struct {
	Format  *model.FormatInfo  `json:"Format"`
	Streams []model.StreamInfo `json:"Streams"`
	Error   *ffprobeError      `json:"Error,omitempty"`
}

// ffprobeError 错误信息
type ffprobeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Parse 是核心解析入口，接收输入字符串，返回标准化的响应结构
// 输入可以是本地文件路径或网络 URL
// 任何错误均被捕获并封装到 Response 中，不会向外抛出 panic
func Parse(input string) *model.Response {
	// 校验输入
	inputType, err := validator.ValidateInput(input)
	if err != nil {
		return errorResponse("invalid input: " + err.Error())
	}

	// 执行 ffprobe
	isNetwork := inputType == validator.TypeRemoteFile
	rawJSON, err := ffmpeg.Probe(ffmpeg.ProbeOptions{
		Input:     input,
		IsNetwork: isNetwork,
	})
	if err != nil {
		return errorResponse("ffprobe execution failed: " + err.Error())
	}

	var probeOut ffprobeOutput
	if err := json.Unmarshal(rawJSON, &probeOut); err != nil {
		return errorResponse("failed to parse ffprobe output: " + err.Error())
	}

	if probeOut.Error != nil && probeOut.Error.Code != 0 {
		return errorResponse(fmt.Sprintf("ffprobe reported error: %s (code: %d)", probeOut.Error.Message, probeOut.Error.Code))
	}

	resp := &model.Response{
		Success: true,
		Format:  probeOut.Format,
		Streams: probeOut.Streams,
	}

	return resp
}

// errorResponse 构造错误响应，统一封装为 JSON 格式
func errorResponse(errMsg string) *model.Response {
	return &model.Response{
		Success: false,
		Error:   errMsg,
	}
}

// GetFFprobeVersion 获取 ffprobe 版本信息
// 若 ffprobe 不可用，返回错误
func GetFFprobeVersion() (string, error) {
	return ffmpeg.GetVersion()
}

// IsFFprobeAvailable 检查 ffprobe 是否可用
func IsFFprobeAvailable() bool {
	return ffmpeg.IsAvailable()
}
