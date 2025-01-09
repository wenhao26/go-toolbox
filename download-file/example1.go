package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v2"
)

const (
	concurrentDownloads = 6                 // 并发下载数，设置为4，表示同时下载4个分块
	timeout             = 300 * time.Second // 每个HTTP请求的超时时间
)

func main() {
	// 下载目标地址和保存路径
	url := "http://prod.publicdata.landregistry.gov.uk.s3-website-eu-west-1.amazonaws.com/pp-complete.csv"
	dest := "F:\\test_files\\download\\pp-complete.csv"

	startTime := time.Now()

	// 开始下载文件
	err := DownloadFile(url, dest)
	if err != nil {
		// 如果下载失败，打印错误信息
		fmt.Printf("下载失败: %v\n", err)
	} else {
		// 如果下载成功，打印成功信息
		fmt.Println("下载已成功完成")
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("下载所耗时为 %s\n", elapsedTime)
}

// DownloadFile 下载文件逻辑
func DownloadFile(url, dest string) error {
	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Minute, // 响应头超时时间
		},
	}

	// 通过HEAD请求获取文件的总大小
	resp, err := client.Head(url)
	if err != nil {
		return fmt.Errorf("无法获取文件大小: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器响应码: %d", resp.StatusCode)
	}

	// 获取文件大小
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return errors.New("无法确定文件大小")
	}

	// 将文件大小从字符串转换为整数
	size, err := strconv.Atoi(contentLength)
	if err != nil {
		return fmt.Errorf("内容长度无效: %w", err)
	}

	fmt.Printf("开始下载: %s (size: %d bytes)\n", url, size)

	// 初始化进度条，用于实时显示下载进度
	progressBar := progressbar.NewOptions(size,
		progressbar.OptionSetDescription("Downloading..."), // 进度条描述
		progressbar.OptionSetWriter(os.Stdout),             // 输出位置为终端
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),                  // 进度条宽度
		progressbar.OptionThrottle(65*time.Millisecond), // 更新频率
		progressbar.OptionShowCount(),                   // 显示已完成块的计数
		progressbar.OptionClearOnFinish(),               // 下载完成后清除进度条
	)

	// 计算每个分块的大小
	partSize := size / concurrentDownloads

	// 保存临时文件的名称
	tmpFiles := make([]string, concurrentDownloads)

	// 用于管理并发任务的同步
	var wg sync.WaitGroup

	// 用于记录并发任务中的错误
	var mu sync.Mutex
	var downloadErr error

	// 启动多个goroutines进行分块下载
	for i := 0; i < concurrentDownloads; i++ {
		wg.Add(1) // 增加等待组计数

		// 计算当前分块的起始和结束字节
		start := i * partSize
		end := start + partSize - 1
		if i == concurrentDownloads-1 {
			// 最后一块下载到文件结束
			end = size - 1
		}

		// 临时文件名
		tmpFile := fmt.Sprintf("part-%d.tmp", i)
		tmpFiles[i] = tmpFile

		// 启动goroutine下载分块
		go func(start, end int, tmpFile string) {
			defer wg.Done() // 下载完成后减少等待组计数
			err := downloadPart(url, start, end, tmpFile, progressBar)
			if err != nil {
				// 如果发生错误，记录错误信息
				mu.Lock()
				downloadErr = err
				mu.Unlock()
			}
		}(start, end, tmpFile)
	}

	// 等待所有goroutines完成
	wg.Wait()

	// 如果有任何分块下载失败，返回错误
	if downloadErr != nil {
		return downloadErr
	}

	// 合并所有分块文件到最终文件
	fmt.Println("开始合并临时文件...")
	err = mergeFiles(tmpFiles, dest)
	if err != nil {
		return fmt.Errorf("合并文件失败: %w", err)
	}

	// 删除临时文件
	fmt.Println("合并成功，删除临时文件...")
	for _, tmpFile := range tmpFiles {
		err := os.Remove(tmpFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

// downloadPart 下载文件的一个分块
func downloadPart(url string, start, end int, dest string, progressBar *progressbar.ProgressBar) error {
	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Minute, // 响应头超时时间
		},
	}

	// 创建HTTP请求，并设置Range头部用于分块下载
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	// 执行HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("服务器未返回部分内容: %d", resp.StatusCode)
	}

	// 打开临时文件用于写入
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将响应体写入文件，并更新进度条
	_, err = io.Copy(io.MultiWriter(out, progressBar), resp.Body)
	return err
}

// mergeFiles 合并所有临时分块文件到最终文件
func mergeFiles(files []string, dest string) error {
	// 打开目标文件用于写入
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// 遍历所有临时文件
	for _, file := range files {
		// 打开临时文件
		in, err := os.Open(file)
		if err != nil {
			return err
		}

		// 将临时文件的内容复制到目标文件
		_, err = io.Copy(out, in)
		in.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
