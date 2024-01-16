package main

import (
	"fmt"
	"strconv"
)

func removeTrailingZeros(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 32) // 将float64转换为字符串形式，保留两位有效数字
}

func main() {
	r1 := float64(1) / float64(2)
	r11 := removeTrailingZeros(r1*100)
	fmt.Printf("%s%s\n", r11, "%")

	r2 := float64(1) / float64(3)
	r22 := removeTrailingZeros(r2*100)
	fmt.Printf("%s%s\n", r22, "%")

}
