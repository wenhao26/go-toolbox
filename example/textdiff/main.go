package main

import (
	"fmt"
	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	t1 := "go-diff 不仅能够简洁地输出字符串对比结果"
	t2 := "go-diffEE 不仅能faf地输出字符串。"

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(t1, t2, false)

	total := len(diffs)
	if total > 0 {
		fmt.Println(total)

		sameCount, different := 0, 0
		for _, diff := range diffs {
			if diff.Type == 0 { // Equal
				sameCount++
			} else {
				different++
			}
			fmt.Println(diff.Type, diff.Text)
		}

		fmt.Println(sameCount, different)
		ret, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sameCount)/float64(total)), 64)
		fmt.Println("相似度为：", (ret*100), "%")
	}

}
