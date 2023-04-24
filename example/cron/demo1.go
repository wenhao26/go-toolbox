package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	// 秒(* / , -)(0-59) 分(* / , -)(0-59) 时(* / , -)(0-23) 日(* / , - ?)(1-31) 月(* / , -)(1-12或者JAN-DEC) 星期(* / , - ?)(0-6或者SUN–SAT)
	/*
	星号（*）
	星号表示匹配该字段的所有值，如在上面表达式的天位置中使用星号，就表示每天。
	斜线（/）
	斜杠用于描述范围的增量，比如'3-59/15'这个表达式在表示从现在的第三分钟开始和往后的每15分钟，到第59分钟为止。表现形式为"* \ / ..."，等同于"N-MAX / m"，即在该字段范围内的增量。即从N开始，使用增量 m 直到 MAX 结束，它没有重复
	逗号（,）
	逗号用于分隔列表中的项，比如，在上表的'星期几'中使用 "MON,WED,FRI" 表示星期一、星期三和星期五
	连字符（-）
	连字符用于定义范围。例如，9-17表示包括上午9点至下午5点在内的每小时
	问号 （?）
	表示不指定值，可以来代替 *
	*/
	_ = c.AddFunc("0-10/2 * * * * *", func() {
		fmt.Println(123)
	})
	c.Start()

	/*ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch*/

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
