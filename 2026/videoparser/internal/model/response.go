package model

// Response 视频解析响应结果 JSON 结构
type Response struct {
	Success bool         `json:"success"`
	Format  *FormatInfo  `json:"format,omitempty"`
	Streams []StreamInfo `json:"streams,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// FormatInfo 存储容器级别的元数据
type FormatInfo struct {
	Filename   string            `json:"filename"`
	NbStreams  int               `json:"nb_streams"`
	FormatName string            `json:"format_name"`
	Duration   string            `json:"duration,omitempty"`
	Size       string            `json:"size,omitempty"`
	BitRate    string            `json:"bit_rate,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
}

// StreamInfo 存储单元流媒体（视频/音频/字母）的详细参数
type StreamInfo struct {
	Index        int               `json:"index"`
	CodecType    string            `json:"codec_type"`
	CodecName    string            `json:"codec_name"`
	ProfileName  string            `json:"profile_name,omitempty"`
	Width        int               `json:"width,omitempty"`
	Height       int               `json:"height,omitempty"`
	PixFmt       string            `json:"pix_fmt,omitempty"`
	SampleRate   string            `json:"sample_rate,omitempty"`
	Channels     int               `json:"channels,omitempty"`
	BitRate      string            `json:"bit_rate,omitempty"`
	Duration     string            `json:"duration,omitempty"`
	RFrameRate   string            `json:"r_frame_rate,omitempty"`
	AvgFrameRate string            `json:"avg_frame_rate,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}
