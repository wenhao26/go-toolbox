package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// ImageTask 图片处理任务
type ImageTask struct {
	ID        int    // 任务唯一标识
	ImagePath string // 图片路径
	Operation string // 行为。`resize`、`watermark`
}

// ProcessedResult 图片处理结果
type ProcessedResult struct {
	TaskID int
	Status string // 处理状态。`success`、`failed`、`cancelled`
	Msg    string // 详细信息
}

const (
	taskCount     = 10 // 任务数量
	workerCount   = 3  // 并发处理的协程数量
	channelBuffer = 5  // 任务信道的缓冲大小
)

// taskProducer 模拟生成图片处理任务，并发送到任务信道
func taskProducer(taskChan chan<- ImageTask, done chan<- bool, taskCount int) {
	fmt.Println("生产者已启动，开始生成图片处理任务...")

	for i := 0; i < taskCount; i++ {
		task := ImageTask{
			ID:        i,
			ImagePath: "path/images/img_" + strconv.Itoa(i) + ".jpg",
			Operation: "resize",
		}

		fmt.Printf("生产者：创建任务 #%d，发送到队列...\n", task.ID)

		taskChan <- task // 将任务发送到任务信道
		time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
	}

	close(taskChan)
	done <- true // 通知主协程生产者已完成任务的创建

	fmt.Println("生产者：所有任务发送完成。")
}

// imageProcessor 模拟图片处理过程
func imageProcessor(id int, taskChan <-chan ImageTask, resultChan chan<- ProcessedResult, cancelChan <-chan struct{}) {
	fmt.Printf("处理者已启动，准备处理任务 #%d...\n", id)

	for {
		select {
		case task, ok := <-taskChan: // 尝试从任务信道接收任务
			if !ok {
				// 任务信道已关闭且无更多数据，说明所有任务已分配完毕
				fmt.Printf("处理者 #%d：任务队列为空，正在退出。\n", id)
				return // 协程退出
			}
			fmt.Printf("处理者 #%d：开始处理任务 #%d...\n", id, task.ID)

			// 模拟处理时间，随机决定是否超时
			processingTime := time.Duration(rand.Intn(4000)+1000) * time.Millisecond // 1秒到5秒
			taskTimeout := time.After(3 * time.Second)                               // 任务最大处理时间3秒

			select {
			case <-time.After(processingTime): // 模拟实际处理时间
				// 任务在规定时间内完成
				fmt.Printf("✅ 处理者 #%d：完成任务 #%d，耗时 %v。\n", id, task.ID, processingTime)
				resultChan <- ProcessedResult{TaskID: task.ID, Status: "success", Msg: "Image processed successfully."}
			case <-taskTimeout: // 任务超时
				fmt.Printf("❌ 处理者 #%d：任务 #%d 超时，耗时 %v，正在取消。\n", id, task.ID, processingTime)
				resultChan <- ProcessedResult{TaskID: task.ID, Status: "cancelled", Msg: "Task processing timed out."}
			case <-cancelChan: // 如果外部发来取消信号
				fmt.Printf("❌ 处理者 #%d：收到取消信号，停止当前任务 #%d 并退出。\n", id, task.ID)
				return // 协程立即退出
			}

		case <-cancelChan: // 如果外部发来取消信号
			fmt.Printf("❌ 处理者 #%d：收到全局取消信号，正在退出。\n", id)
			return // 协程立即退出
		}
	}
}

// resultCollector 收集处理结果并打印
func resultCollector(resultChan <-chan ProcessedResult, done chan<- bool, expectedResults int) {
	fmt.Println("结果收集者启动：等待收集处理结果...")

	collectedCount := 0
	for result := range resultChan {
		fmt.Printf("结果收集者：收到任务 #%d 结果 -> 状态: %s, 消息: %s\n", result.TaskID, result.Status, result.Msg)
		collectedCount++
		if collectedCount >= expectedResults {
			break // 达到预期结果数量，提前退出循环
		}
	}
	done <- true // 通知主协程结果收集者已完成
	fmt.Println("结果收集者完成：已收集所有预期结果。")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 创建信道
	taskChannel := make(chan ImageTask, channelBuffer)     // 任务信道，带缓冲
	resultChannel := make(chan ProcessedResult, taskCount) // 结果信道，缓冲大小等于任务总数
	producerDone := make(chan bool)                        // 生产者完成信号信道
	collectorDone := make(chan bool)                       // 结果收集者完成信号信道
	globalCancel := make(chan struct{})                    // 全局取消信号信道

	fmt.Println("--- 程序开始：图片处理系统 ---")

	// 启动生产者协程
	go taskProducer(taskChannel, producerDone, taskCount)

	// 启动多个处理者协程
	for i := 0; i < workerCount; i++ {
		go imageProcessor(i, taskChannel, resultChannel, globalCancel)
	}

	// 启动结果收集者协程
	go resultCollector(resultChannel, collectorDone, taskCount)

	// 等待生产者完成任务发送
	<-producerDone
	fmt.Println("\n--- 生产者已完成任务发送，等待所有任务处理完毕或超时 ---")

	// 给处理者一些时间完成任务
	// 实际应用中，这里可能需要更复杂的逻辑来判断所有任务是否真的处理完毕
	// 例如：等待所有 worker 协程都退出，或者等待收集到所有预期的结果
	// 这里我们通过等待结果收集者完成来判断
	<-collectorDone

	// 关闭结果信道，通知所有处理者可以安全退出了
	// 注意：这里需要确保所有处理者都已从 taskChannel 读取完毕或因超时退出
	// 如果处理者仍在等待 taskChannel，它不会收到关闭信号。
	// 在本例中，因为 taskChannel 已在生产者中关闭，处理者会因此退出。
	close(resultChannel) // 关闭结果信道，不再接收结果

	// 如果需要提前停止所有处理者，可以发送取消信号
	// close(globalCancel) // 在这里发送，所有处理者会立即退出

	fmt.Println("\n--- 所有任务处理及结果收集完成，程序即将退出 ---")

	// 最终阻止主协程退出，直到所有协程都真正完成
	// 因为我们上面已经等待了 collectorDone，主协程知道结果收集完成了
	// 但仍需要确保所有 worker 协程也已退出，这里用一个短暂的延时来代替更复杂的同步机制
	time.Sleep(2 * time.Second)
	fmt.Println("--- 程序结束 ---")
}
