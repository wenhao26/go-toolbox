//go:build !windows

package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

// registerPlatformSignals 在 Linux/macOS 环境下挂载高级特异性信号
func (osm *OmniShutdownManager) registerPlatformSignals(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
}

// isReloadSignal 专门在 Unix/Linux 下判断是否为重载信号
func (osm *OmniShutdownManager) isReloadSignal(sig os.Signal) bool {
	return sig == syscall.SIGUSR1 || sig == syscall.SIGUSR2
}
