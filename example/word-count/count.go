package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

func TotalWords(s string) int {
	wordCount := 0

	plainWords := strings.Fields(s)
	for _, word := range plainWords {
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			wordCount++
		} else {
			wordCount += runeCount
		}
	}

	return wordCount
}

// 统计多语言字符串的字数函数
func CountWords(input string) int {
	// 定义正则表达式
	// [\p{Han}] 匹配单个中文字符
	// [\p{Hangul}] 匹配单个韩文字符
	// [a-zA-Z0-9]+ 匹配英文单词和数字
	// [\p{Thai}\p{L}]+ 匹配泰语和其他语言的单词
	re := regexp.MustCompile(`[\p{Han}]|[\p{Hangul}]|[a-zA-Z0-9]+|[\p{Thai}\p{L}]+`)

	// 使用正则表达式提取符合条件的片段
	matches := re.FindAllString(input, -1)

	// 返回匹配的数量作为字数
	return len(matches)
}

func main() {
	// text := "写这样一个工具的原因是部分 Light weight 的 web framework 并没有内置自动重新加载的功能"
	//text := "hello,web framework"
	//count := TotalWords(text)
	//fmt.Println(count)

	text := "你好, world! Selamat pagi, สวัสดี 123 😊. Chào bạn! "
	wordCount := CountWords(text)
	fmt.Printf("字数统计结果: %d\n", wordCount)
}
