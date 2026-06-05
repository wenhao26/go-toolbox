package shutdown

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Hook 定义组件退出时需要执行的清理行为契约
type Hook func(ctx context.Context) error

// ReloadHook 定义配置热重载（不 kill 进程）行为契约
type ReloadHook func() error

// Option 定义配置函数类型
type Option func(manager *OmniShutdownManager)

// WithTimeout 设置优雅退出的硬超时时间
// 如果不设置或者 ≤0，则自动切换为“无限制等待”模式
func WithTimeout(timeout time.Duration) Option {
	return func(osm *OmniShutdownManager) {
		osm.timeout = timeout
	}
}

// OmniShutdownManager 负责全场景的生命周期编排
// - 外部信号
// - 内部 Panic
// - 配置重载
type OmniShutdownManager struct {
	timeout     time.Duration
	mu          sync.Mutex
	hooks       []Hook
	reloadHooks []ReloadHook
	isClosed    int32 // 原子操作标记，保障全局退出流程只被触发一次
}

// New 创建一个开箱即用的全场景优雅退出管理器
func New(opts ...Option) *OmniShutdownManager {
	osm := &OmniShutdownManager{
		hooks:       make([]Hook, 0),
		reloadHooks: make([]ReloadHook, 0),
	}

	for _, opt := range opts {
		opt(osm)
	}

	return osm
}

// RegisterShutdown 注册停机清理钩子
// 底层采用后进后出（LIFO）栈设计。先注册的组件，最后关闭（如：先关闭流量入口，后释放底层数据库）
func (osm *OmniShutdownManager) RegisterShutdown(hook Hook) {
	if hook == nil {
		return
	}

	osm.mu.Lock()
	defer osm.mu.Unlock()

	osm.hooks = append(osm.hooks, hook)
}

// RegisterReload 注册热重载钩子
// 当收到动态配置刷新信号时触发
func (osm *OmniShutdownManager) RegisterReload(hook ReloadHook) {
	if hook == nil {
		return
	}

	osm.mu.Lock()
	defer osm.mu.Unlock()

	osm.reloadHooks = append(osm.reloadHooks, hook)
}

// SystemPanicGuard 提供全局的 Panic 捕获哨兵
// 常驻进程的主写成或关键异步工作协程应该直接 `defer osm.SystemPanicGuard()` 挂载
func (osm *OmniShutdownManager) SystemPanicGuard() {
	if r := recover(); r != nil {
		log.Printf("[OmniEngine] 🚨 致命异常: 检测到未捕获的运行时 Panic: %v\n堆栈轨迹:\n%s", r, string(debug.Stack()))
		log.Println("[OmniEngine] 正在将内部 Panic 转化为优雅退出行为，尝试挽救在途核心数据...")
		osm.triggerShutdown(fmt.Sprintf("Internal-Panic: %v", r))
	}
}

// WaitListen 核心阻塞方法：拉起全场景信号监控矩阵
// 返回值：建议的进程退出状态码（0 表示圆满成功，1 表示异常或被强制击杀）
func (osm *OmniShutdownManager) WaitListen() int {
	// 针对 SIGPIPE（管道破裂）进行忽略，防止常驻服务因为一次突发的底层网络 Socket 写入失败而直接猝死
	osm.ignoreSignal(syscall.SIGPIPE)

	// 创建带缓冲的信号通道，确保高并发系统突发多组信号时不阻塞
	sigChan := make(chan os.Signal, 10)

	// 调用跨平台多信号挂载函数
	osm.notifySignals(sigChan)
	log.Println("[OmniEngine] 🚀 全场景动态场景监控矩阵已成功挂载，常驻进程监听中...")

	for {
		sig := <-sigChan
		if sig == nil {
			return 1
		}

		// 执行平台特异性判断，检查这个信号是不是该平台的“热重载信号”
		if osm.isReloadSignal(sig) {
			log.Printf("[OmniEngine] 🔍 场景监控: 检测到运维热重载配置指令 (%v)。开始执行配置无感刷新...", sig)
			osm.executeReload()
			continue // 继续循环监听，不退出进程
		}

		switch sig {
		case syscall.SIGHUP:
			log.Println("[OmniEngine] 🔍 场景监控: 检测到远程终端/控制台已断开连接 (SIGHUP)。")
			return osm.handleShutdownFlow("SIGHUP", sigChan)

		case syscall.SIGINT:
			log.Println("[OmniEngine] 🔍 场景监控: 检测到键盘人工中断信号 (SIGINT / Ctrl+C)。")
			return osm.handleShutdownFlow("SIGINT", sigChan)

		case syscall.SIGTERM:
			log.Println("[OmniEngine] 🔍 场景监控: 检测到系统销毁指令 (SIGTERM / K8s Pod 滚动更新)。")
			return osm.handleShutdownFlow("SIGTERM", sigChan)

		default:
			log.Printf("[OmniEngine] 🔍 场景监控: 捕获到其他未定义系统信号 [%v]，跳过处理。\n", sig)
		}
	}
}

// triggerShutdown 内部组件主动诱发优雅退出
func (osm *OmniShutdownManager) triggerShutdown(name string) {
	if !atomic.CompareAndSwapInt32(&osm.isClosed, 0, 1) {
		return
	}

	ctx, cancel := osm.buildContext()
	defer cancel()

	_ = osm.executeHooks(ctx)
	os.Exit(1)
}

// handleShutdownFlow 核心状态机：收拢处理外部信号的退出逻辑
func (osm *OmniShutdownManager) handleShutdownFlow(reason string, sigChan chan os.Signal) int {
	if !atomic.CompareAndSwapInt32(&osm.isClosed, 0, 1) {
		return 0
	}

	ctx, cancel := osm.buildContext()
	defer cancel()

	// 异步驱动清理链条
	done := make(chan error, 1)
	go func() {
		done <- osm.executeHooks(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("[OmniEngine] ❌ 优雅停机执行期间发生错误: %v\n", err)
			return 1
		}
		log.Printf("[OmniEngine] 🎉 场景 [%s] 优雅停机圆满结束。\n", reason)
		return 0

	case sig2 := <-sigChan:
		// 工业级双击强杀守卫：如果用户或运维在等待期间再次发送 Ctrl+C 或 kill 信号，视为强行终止
		if sig2 == syscall.SIGINT || sig2 == syscall.SIGTERM {
			log.Printf("[OmniEngine] ⚡ 强杀监控: 捕获到二次强杀指令: %v。进程立即闪退！\n", sig2)
			return 1
		}
		return 1
	}
}

// buildContext 构建上下文
func (osm *OmniShutdownManager) buildContext() (context.Context, context.CancelFunc) {
	if osm.timeout > 0 {
		return context.WithTimeout(context.Background(), osm.timeout)
	}
	return context.WithCancel(context.Background())
}

// executeHooks 执行钩子
func (osm *OmniShutdownManager) executeHooks(ctx context.Context) error {
	osm.mu.Lock()
	hooksCopy := make([]Hook, len(osm.hooks))
	copy(hooksCopy, osm.hooks)
	osm.mu.Unlock()

	// 逆序安全释放资源
	for i := len(hooksCopy) - 1; i >= 0; i-- {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := hooksCopy[i](ctx); err != nil {
			log.Printf("[OmniEngine] 资源释放失败 [Index: %d]: %v\n", i, err)
		}
	}

	return nil
}

func (osm *OmniShutdownManager) executeReload() {
	osm.mu.Lock()
	reloads := make([]ReloadHook, len(osm.reloadHooks))
	copy(reloads, osm.reloadHooks)
	osm.mu.Unlock()

	for _, reload := range reloads {
		if err := reload(); err != nil {
			log.Printf("[OmniEngine] 配置热重载失败: %v\n", err)
		}
	}
}

func (osm *OmniShutdownManager) ignoreSignal(sig os.Signal) {
	// 防御捕获 Windows 平台下可能发生的未定义引发的 panic
	defer func() { _ = recover() }()
	signal.Ignore(sig)
}

func (osm *OmniShutdownManager) notifySignals(sigChan chan os.Signal) {
	// 全平台通用的核心必选信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// 动态挂载平台特异性信号
	osm.registerPlatformSignals(sigChan)
}
