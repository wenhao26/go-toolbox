package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"toolbox/2026/omni-shutdown/pkg/shutdown"
)

// AIContentWorker 模拟处理硬核慢长任务的常驻消费者
type AIContentWorker struct {
	doneChan chan struct{}
}

func (w *AIContentWorker) StartConsumeLoop(ctx context.Context, engine *shutdown.OmniShutdownManager) {
	defer close(w.doneChan)
	// 规范：每一个在常驻进程里拉起的独立异步底层协程，内部都必须防御性挂载 PanicGuard
	defer engine.SystemPanicGuard()

	for taskID := 1; taskID <= 3; taskID++ {
		select {
		case <-ctx.Done():
			log.Println("[AI-Worker] 流量通道已关闭，不再拉取新小说任务，退出处理死循环。")
			return
		default:
		}

		w.executeSevereTask(ctx, taskID)

		// 模拟线上偶发的运行时致命异常场景：在处理完第 2 个任务后，系统意外发生越界 Panic
		if taskID == 2 {
			log.Println("[AI-Worker] ⚠️ 触发黑天鹅突发事件：任务协程内部发生未知异常崩溃...")
			var arr []string
			_ = arr[1024] // 人为诱发运行时越界 Panic，由 SystemPanicGuard 捕获并强制带离，启动优雅退场
		}
	}
}

func (w *AIContentWorker) executeSevereTask(ctx context.Context, id int) {
	log.Printf("[AI-Worker] 🚀 正在处理硬核长任务 [%d]，调用 Azure AI 进行精细化文本上下文推理...\n", id)

	// 业务内部实现精细的检查点契约（Checkpoint）
	for step := 1; step <= 3; step++ {
		time.Sleep(1 * time.Second) // 模拟复杂的网络轮询
		fmt.Printf("   -> 任务 [%d] 离线生成进度: Step %d/3 完成\n", id, step)
	}
	log.Printf("[AI-Worker]  长任务 [%d] 离线生成质量达标，落库成功。\n", id)
}

func main() {
	// 场景二：常驻 AI 自动化长任务生成（无限制等待模式 + 内部 Panic 转化 + 配置重载）
	// 适用于单次物理任务耗时较长（如 15~20 分钟的 AI 长文本小说自动生成、海量数据水平同步切片）。
	// 其特征是不能中途强杀以防丢失巨额算力，但必须防范内部协程 Panic 并支持在线配置文件刷新。
	log.Println("[AI-Engine] 正在拉起常驻大模型自动化生成生产引擎...")

	// 1. 初始化引擎：不传 WithTimeout 选项，自动开启【无限制等待模式】
	engine := shutdown.New()

	// 2. 挂载主协程全局的 Panic 防护盾
	defer engine.SystemPanicGuard()

	// 3. 注册大厂运维最常用的“配置无感热重载” Hook (Linux 下通过命令 kill -USR1 <pid> 触发)
	engine.RegisterReload(func() error {
		log.Println("[Reload] 📡 刷新信号到达：成功重新拉取配置，各账号 API Rate-Limit 计数器复位。")
		return nil
	})

	// 4. 创建常驻上下文，用于在优雅停机时向下游在途业务下发“不再接收新任务”的契约宣告
	pipelineCtx, stopPipeline := context.WithCancel(context.Background())
	defer stopPipeline()

	// 5. 初始化并异步拉起常驻任务消费者
	worker := &AIContentWorker{doneChan: make(chan struct{})}
	go worker.StartConsumeLoop(pipelineCtx, engine)

	// 6. 编排停机序列：步骤 A（率先执行）通知 Worker 拦截新流量，并死等手里已有的慢任务跑完
	engine.RegisterShutdown(func(ctx context.Context) error {
		log.Println("[Shutdown-Hook] 收到指令！正在撤回生产线接收器，拒绝承接新业务...")
		stopPipeline() // 触发 context 信号

		<-worker.doneChan // 阻塞，无限期死等 Worker 报告手里当前的在途长任务完美收尾
		return nil
	})

	// 7. 编排停机序列：步骤 B（最后执行）关闭底座设施
	engine.RegisterShutdown(func(ctx context.Context) error {
		log.Println("[Shutdown-Hook] 正在注销云端大模型长连接握手及文件持久化句柄...")
		return nil
	})

	// 8. 移交控制权。如果卡死，支持运维现场“再次敲击 Ctrl+C”触发双击强杀硬退
	code := engine.WaitListen()
	os.Exit(code)
}
