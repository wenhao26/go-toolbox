package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	url := "https://seopic.699pic.com/photo/50061/8976.jpg_wh1200.jpg"
	output := "image.jpg"

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发送 GET 请求获取响应
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// 创建输出文件
	outFile, err := os.Create(output)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer outFile.Close()

	// 获取文件大小
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// 创建进度条
	progress := NewProgressBar(size)

	// 创建多个管道
	writer := io.MultiWriter(outFile, progress)

	// 复制响应主体到文件
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Download complete")
}

// 进度条结构体
type ProgressBar struct {
	Total   int
	Current int
}

// 创建新的进度条
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total:   total,
		Current: 0,
	}
}

// 实现 io.Writer 接口
func (p *ProgressBar) Write(b []byte) (n int, err error) {
	n = len(b)
	p.Current += n
	p.Print()
	return
}

// 打印进度条
func (p *ProgressBar) Print() {
	if p.Total == 0 {
		return
	}
	percent := float64(p.Current) / float64(p.Total) * 100
	fmt.Printf("\r%d/%d bytes (%.2f%%)", p.Current, p.Total, percent)
}
