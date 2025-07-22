package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// - done chan struct{} 的使用案例

// worker1 基础用法：简单退出信号
func worker1(done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("收到退出信号，停止工作")
			return
		default:
			fmt.Println("运行中...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 组合使用：带超时控制的退出
func worker2(done chan struct{}) (string, error) {
	select {
	case <-time.After(1 * time.Second):
		return "数据内容", nil
	case <-done:
		return "", fmt.Errorf("请求被取消")
	}
}

// 多协程协同退出
func worker3(id int, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-done:
			fmt.Printf("worker %d 收到退出信号\n", id)
			return
		default:
			fmt.Printf("worker %d 正在工作\n", id)
			// time.Sleep(time.Duration(id) * 300 * time.Millisecond)
		}
	}
}

// 与 context 结合使用
func worker4(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s 收到取消信号：%v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s 正在运行\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 高级模式：多级退出控制
func worker5a(done <-chan struct{}) {
	secondaryDone := make(chan struct{})

	go worker5b(secondaryDone)

	for {
		select {
		case <-done:
			fmt.Println("主worker收到退出信号")
			close(secondaryDone) // 通知次级worker退出
			return
		default:
			fmt.Println("主worker工作中...")
			time.Sleep(800 * time.Millisecond)
		}
	}
}

func worker5b(done <-chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("次级worker收到退出信号")
			return
		default:
			fmt.Println("次级worker工作中...")
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan struct{})

	// ==== 示例一 ====
	//go worker1(done)
	//
	//time.Sleep(2 * time.Second)
	//close(done)
	//
	//time.Sleep(500 * time.Millisecond)

	// ==== 示例二 ====
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	close(done)
	//}()
	//
	//result, err := worker2(done)
	//if err != nil {
	//	fmt.Println("错误：", err)
	//	return
	//}
	//fmt.Println("结果：", result)

	// ==== 示例三 ====
	//var wg sync.WaitGroup
	//
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go worker3(i, done, &wg)
	//}
	//
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	close(done)
	//}()
	//wg.Wait()

	// ==== 示例四 ====
	// 创建一个可取消的context
	//ctx, cancel := context.WithCancel(context.Background())
	//
	//go worker4(ctx, "worker1")
	//go worker4(ctx, "worker2")
	//
	//time.Sleep(2 * time.Second)
	//cancel() // 取消所有worker

	// ==== 示例五 ====
	go worker5a(done)
	time.Sleep(2 * time.Second)
	close(done)

	fmt.Println("程序结束", done)
}
