// 使用ants协程池来并发下载图片
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const saveDir = "F:\\test_files\\download\\images"

// downloadImage 下载图片文件（普通模式）
//func downloadImage(imgUrl string) error {
//	// 创建HTTP客户端
//	client := &http.Client{
//		Timeout: 10 * time.Second, // 设置超时时间
//	}
//
//	// 发送HTTP GET请求
//	resp, err := client.Get(imgUrl)
//	if err != nil {
//		return fmt.Errorf("请求失败: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// 检查响应状态码
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
//	}
//
//	// 处理 URL，提取文件名（去掉查询参数部分）
//	parsedURL := imgUrl
//	if idx := strings.Index(parsedURL, "?"); idx != -1 {
//		parsedURL = parsedURL[:idx] // 去掉 URL 中的查询参数
//	}
//
//	// 从URL中提取文件名
//	_, filename := path.Split(parsedURL)
//	if filename == "" {
//		// 如果没有文件名，可以使用默认的名称
//		filename = "image.jpg"
//	}
//
//	savePath := fmt.Sprintf("%s\\%s", saveDir, filename)
//
//	// 创建保存图片的文件
//	file, err := os.Create(savePath)
//	if err != nil {
//		return fmt.Errorf("创建文件失败: %v", err)
//	}
//	defer file.Close()
//
//	// 将图片内容写入文件
//	_, err = io.Copy(file, resp.Body)
//	if err != nil {
//		return fmt.Errorf("写入文件失败: %v", err)
//	}
//
//	fmt.Printf("图片下载成功，保存路径: %s\n", savePath)
//	return nil
//}

// downloadImage 下载图片文件（协程并发模式）
func downloadImage(imgUrl string, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间
	}

	// 发送HTTP GET请求
	resp, err := client.Get(imgUrl)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 处理 URL，提取文件名（去掉查询参数部分）
	parsedURL := imgUrl
	if idx := strings.Index(parsedURL, "?"); idx != -1 {
		parsedURL = parsedURL[:idx] // 去掉 URL 中的查询参数
	}

	// 从URL中提取文件名
	_, filename := path.Split(parsedURL)
	if filename == "" {
		// 如果没有文件名，可以使用默认的名称
		filename = "image.jpg"
	}

	savePath := fmt.Sprintf("%s\\%s", saveDir, filename)

	// 创建保存图片的文件
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 将图片内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("图片下载成功，保存路径: %s\n", savePath)
	return nil
}

func main() {
	imgUrls := []string{
		"https://images.pexels.com/photos/8952192/pexels-photo-8952192.jpeg?cs=srgb&dl=pexels-ivan-samkov-8952192.jpg&fm=jpg",
		"https://images.pexels.com/photos/14996824/pexels-photo-14996824.jpeg?cs=srgb&dl=pexels-midtrack-14996824.jpg&fm=jpg",
		"https://images.pexels.com/photos/20263436/pexels-photo-20263436.jpeg?cs=srgb&dl=pexels-yuraforrat-20263436.jpg&fm=jpg",
		"https://images.pexels.com/photos/17642974/pexels-photo-17642974.jpeg?cs=srgb&dl=pexels-ozanculha-17642974.jpg&fm=jpg",
		"https://images.pexels.com/photos/4344756/pexels-photo-4344756.jpeg?cs=srgb&dl=pexels-apasaric-4344756.jpg&fm=jpg",
		"https://images.pexels.com/photos/15447422/pexels-photo-15447422.jpeg?cs=srgb&dl=pexels-kyle-miller-169884138-15447422.jpg&fm=jpg",
		"https://images.pexels.com/photos/29990902/pexels-photo-29990902.jpeg?cs=srgb&dl=pexels-tiarrasorte-29990902.jpg&fm=jpg",
		"https://images.pexels.com/photos/28718329/pexels-photo-28718329.jpeg?cs=srgb&dl=pexels-eslames1-28718329.jpg&fm=jpg",
		"https://images.pexels.com/photos/30149969/pexels-photo-30149969.jpeg?cs=srgb&dl=pexels-raymond-petrik-1448389535-30149969.jpg&fm=jpg",
		"https://images.pexels.com/photos/17516413/pexels-photo-17516413.jpeg?cs=srgb&dl=pexels-wojtekpaczes-17516413.jpg&fm=jpg",
	}

	startTime := time.Now()

	// 普通模式
	//for _, imgUrl := range imgUrls {
	//	fmt.Println(imgUrl)
	//	err := downloadImage(imgUrl)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	// 协程并发模式
	var wg sync.WaitGroup

	// 创建协程池
	pool, err := ants.NewPool(4)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	for _, imgUrl := range imgUrls {
		wg.Add(1)

		imgUrlCopy := imgUrl
		_ = pool.Submit(func() {
			err := downloadImage(imgUrlCopy, &wg)
			if err != nil {
				panic(err)
			}
		})
	}

	wg.Wait()

	fmt.Printf("抽奖流程成功完成，耗时: %v\n", time.Since(startTime))
}
