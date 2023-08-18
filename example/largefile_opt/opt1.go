package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var filePath = "F:\\test_files\\test.txt"
var short = "VARCHAR"
var tableValue []int

func readFile(filePath *string) *[]string {
	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	ch := make(chan int)
	log.Println("加载文件...")
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("")
				log.Println("文件加载完成")
				return
			default:
				time.Sleep(time.Second)
				fmt.Print(".")
			}
		}
	}()

	var strArr []string
	buf := make([]byte, 2048)
	for i := 0; ; i++ {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		strArr = append(strArr, string(buf[:n]))
	}
	ch <- 1
	return &strArr
}

func arithmeticKmp(haystack string, needle string) int {
	index := -1
	if tableValue == nil {
		tableValue = getPartialMatchTable(needle)
	}
	i, j := 0, 0
	for ; i < len(haystack) && j < len(needle); {
		if haystack[i] == needle[j] {
			if j == 0 {
				index = i
			} // 记录第一个匹配字符的索引
			j++
			i++
		} else {
			if j == 0 {
				i = i + j + 1 - tableValue[j] // 移动位数 =已匹配的字符数 - 对应的部分匹配值
			} else {
				i = index + j - tableValue[j-1] // 如果已匹配的字符数不为零，则重新定义i迭代
			}
			j = 0 // 将已匹配迭代置为0
			index = -1
		}
	}
	return index
}

func getPartialMatchTable(str string) []int {
	var left, right []string // 前缀、后缀
	n := len([]rune(str))
	var result = make([]int, n)
	for i := 0; i < n; i++ {
		left = make([]string, i)  // 实例化前缀容器
		right = make([]string, i) // 实例化后缀容器
		// 前缀
		for j := 0; j < i; j++ {
			if j == 0 {
				left[j] = string(str[j])
			} else {
				left[j] = left[j-1] + string(str[j])
			}
		}
		// 后缀
		for k := i; k > 0; k-- {
			if k == i {
				right[k-1] = string(str[k])
			} else {
				right[k-1] = string(str[k]) + right[k]
			}
		}
		// 找到前缀和后缀中相同的项，长度即为相等项的长度（相等项应该只有一项）
		num := len(left) - 1
		for m := 0; m < len(left); m++ {
			if right[num] == left[m] {
				result[i] = len(left[m])
			}
			num--
		}
	}
	return result
}

func getCount(long string, short string, ch chan int) {
	if len(short) > len(long) {
		ch <- 0
	}

	num := 0
	for ; len(long) > len(short); {
		index := arithmeticKmp(long, short)
		if index == -1 {
			break
		}
		num++

		if index+len(short) > len(long) {
			break
		}
		long = long[index+len(short):]
	}
	ch <- num
}

func getResult(ch chan int, request chan int32) {
	// time.Sleep(time.Second)
	var num int32 = 0

	for {
		select {
		case t := <-ch:
			atomic.AddInt32(&num, int32(t))
		default:
			request <- atomic.LoadInt32(&num)
		}
	}
}

func main() {
	log.Println("开始计时")
	// 读取文件
	longs := readFile(&filePath)

	//log.Println("开始计时")
	t1 := time.Now()
	ch := make(chan int, 100)
	request := make(chan int32)

	for _, long := range *longs {
		go getCount(long, short, ch)
	}

	go getResult(ch, request)
	log.Println(<-request)

	elapsed := time.Since(t1)
	fmt.Println("程序消耗时间: ", elapsed)
}
