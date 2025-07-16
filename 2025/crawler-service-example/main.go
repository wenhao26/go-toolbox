package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// - 模拟了一个简化的网络爬虫系统

type URLData struct {
	URL     string
	Content string
	Error   error
}

// fetchURL 负责抓取URL内容
func fetchURL(url string, dataChan chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("开始抓取：%s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		dataChan <- URLData{
			URL:   url,
			Error: fmt.Errorf("抓取失败: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dataChan <- URLData{
			URL:   "",
			Error: fmt.Errorf("读取响应失败: %v", err),
		}
		return
	}

	dataChan <- URLData{
		URL:     url,
		Content: string(body),
	}
	fmt.Printf("完成抓取：%s\n", url)
}

// parseContent 解析抓取内容
func parseContent(dataChan <-chan URLData, resultChan chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done()

	for urlData := range dataChan {
		if urlData.Error != nil {
			resultChan <- urlData
			continue
		}

		fmt.Printf("开始解析: %s\n", urlData.URL)

		// 模拟解析操作，比如取标题、链接等耗时操作
		time.Sleep(time.Duration(200+int(time.Now().UnixNano()%300)) * time.Millisecond)

		urlData.Content = "已解除内容长度：" + strconv.Itoa(len(urlData.Content))
		resultChan <- urlData

		fmt.Printf("完成解析: %s\n", urlData.URL)
	}
}

func main() {
	urls := []string{
		"http://www.google.com",
		"http://httpbin.org/delay/2",        // 模拟一个慢响应URL
		"http://httpbin.org/delay/4",        // 模拟一个慢响应URL
		"http://nonexistent-domain-xyz.com", // 模拟一个会导致错误URL
		"https://www.zhihu.com/",
	}

	dataChan := make(chan URLData, len(urls))
	resultChan := make(chan URLData, len(urls))

	var wg sync.WaitGroup

	fmt.Println("---- 爬虫服务启动 ----")

	// 启动协程处理URL抓取
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, dataChan, &wg)
	}

	// 启动协程解析URL内容
	wg.Add(1)
	go parseContent(dataChan, resultChan, &wg)

	go func() {
		wg.Wait()
		close(dataChan)
		fmt.Println("所有抓取和解析协程已完成，dataChan 已关闭")
	}()

	// 结果处理
	var collectedResults []URLData
	var resWg sync.WaitGroup // 另一个 WaitGroup 来等待结果收集协程完成

	resWg.Add(1)
	go func() {
		defer resWg.Done()
		for result := range resultChan {
			collectedResults = append(collectedResults, result)
		}
	}()

	// 在主协程中，等待 dataChan 的所有写入和读取都完成
	time.Sleep(time.Second * 5)
	close(resultChan)
	resWg.Wait()

	// 使用 `select{}` 阻塞主协程
	//select {} // 如果这里不需要打印汇总结果，或者需要服务一直运行，可以用它来阻塞

	fmt.Println("---- 最终结果汇总 ----")
	for _, data := range collectedResults {
		if data.Error != nil {
			fmt.Printf("错误: URL: %s, 错误: %v\n", data.URL, data.Error)
		} else {
			fmt.Printf("成功: URL: %s, 内容: %s\n", data.URL, data.Content)
		}
	}
	fmt.Println("Done")
}
