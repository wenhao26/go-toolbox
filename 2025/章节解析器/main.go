// 将章节标题+内容分离
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Chapter 章节结构体
type Chapter struct {
	Title   string
	Content string
}

func main() {
	// 打开文件
	file, err := os.Open("chapter.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建一个scanner来读取文件
	scanner := bufio.NewScanner(file)

	var chapters []Chapter
	var currentChapter *Chapter

	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否是章节标题
		if strings.HasPrefix(line, "#####") {
			// 如果当前章节不为空，将其添加到章节列表中
			if currentChapter != nil {
				chapters = append(chapters, *currentChapter)
			}

			// 创建新的章节
			title := strings.TrimSpace(strings.TrimPrefix(line, "#####"))
			currentChapter = &Chapter{Title: title}
		} else if currentChapter != nil {
			// 如果不是标题行，且当前章节不为空，则将内容添加到当前章节
			currentChapter.Content += line + "\n"
		}
	}

	// 添加最后一个章节
	if currentChapter != nil {
		chapters = append(chapters, *currentChapter)
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// 输出解析结果
	for _, chapter := range chapters {
		fmt.Printf("Title: %s\n", chapter.Title)
		fmt.Printf("Content:\n%s\n", chapter.Content)
		fmt.Println("-----")
	}
}
