package main

import (
	"fmt"
	"time"
)

func WriteDate(n int, c chan int) {
	for i := 1; i <= n; i++ {
		fmt.Println("协助程写入数据：", i)
		c <- i * n //管道只有1个大小，只能1个1个写入
		//只要编译器发现有读取这个管道的操作这里就会阻塞直到读取后让出位置再写入，没有读管道会报错塞满

	}
	close(c) //当管道完成自己的任务后就应该关闭，否则使用FOR读到一个未关闭的管道也会报错
}

func ReadDate(c chan int, b chan bool) {
	for {
		v, ok := <-c
		if !ok {
			break
		}
		fmt.Println("协程读取数据：", v)
		time.Sleep(time.Second) //休眠1秒测试GO编译器的异步处理机制
	}
	b <- true //读取的任务完成了后给锁主进程的FOT一个通知，可以继续向下执行了
	close(b)  //通知已完成，管道需要关闭
}
func main() {
	var c = make(chan int, 1)  //管道大小设置为1，用于测试go编译器自动阻塞机制
	var b = make(chan bool, 1) //用于锁住主进程直到协程要做的事完成后才返回数据

	go WriteDate(10, c) //启动1个协程写入数据
	go ReadDate(c, b)   //再启动一个协程处理读数据

	for f := 1; ; f++ {
		fmt.Println("使用for无限循环锁定主进程！", f)
		_, ok := <-b
		if !ok {
			fmt.Println("数据读取完毕！解除锁定！", f)
			break
		}
	}
}
