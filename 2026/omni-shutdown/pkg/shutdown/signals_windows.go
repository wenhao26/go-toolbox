//go:build windows

package shutdown

import (
	"os"
)

// registerPlatformSignals 针对 Windows 环境进行平滑降级，不挂载 Unix 独有信号，确保 Windows 10 下编译通畅
func (osm *OmniShutdownManager) registerPlatformSignals(sigChan chan os.Signal) {
	// Windows 平台下只依赖通用的 SIGINT 和 SIGTERM 即可
}

// isReloadSignal Windows 平台直接返回 false，绝不引用不存在的系统常量
func (osm *OmniShutdownManager) isReloadSignal(sig os.Signal) bool {
	return false
}
