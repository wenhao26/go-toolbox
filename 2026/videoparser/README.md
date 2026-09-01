# Video Parser - 高性能视频信息解析工具

基于 Golang + FFmpeg 实现的视频元信息解析 CLI 工具，支持本地文件和网络视频链接。

## 功能特性

- 支持本地视频文件解析（MP4、AVI、MKV、MOV 等常见格式）
- 支持网络视频链接解析（HTTP/HTTPS）
- 输出结构化 JSON 数据，包含容器格式、视频流、音频流等详细信息
- 跨平台支持（Windows、Linux、macOS）
- 高性能解析（基于 ffprobe，不解码视频）

## 安装

### 前置依赖

- **FFmpeg**（包含 ffprobe）
    - Windows: 从 [gyan.dev](https://www.gyan.dev/ffmpeg/builds/) 下载并添加到 PATH
    - Linux: `sudo apt install ffmpeg` 或 `sudo yum install ffmpeg`
    - macOS: `brew install ffmpeg`

验证安装：
```bash
ffprobe -version