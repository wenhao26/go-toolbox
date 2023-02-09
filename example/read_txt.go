package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func streamReadFile(filepath string) int {
	count := 0

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		// line, err := buf.ReadString('\n') // 使用此方式无法读取最后一行
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				fmt.Println("文本已读取完成")
				break
			}
			fmt.Printf("读取文件失败,错误为:%v", err)
			break
		}

		fmt.Println(string(line))

		count += totalWords(string(line))
		fmt.Println(count)
	}

	return count
}

func sliceReadFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := make([]byte, 4096)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			_, _ = fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
			os.Exit(1)
		case nr == 0: // EOF
			fmt.Println("文本已读取完成")
			return
		case nr > 0:
			fmt.Println(string(s[:nr]))
		}
	}
}

func totalWords(s string) int {
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

func main() {
	/*totalWords := streamReadFile("data/novel.txt")
	fmt.Println(totalWords)*/

	sliceReadFile("data/novel.txt")
}
