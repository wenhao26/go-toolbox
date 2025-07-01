package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel 定义日志级别。
type LogLevel int

const (
	DEBUG LogLevel = iota // 调试信息
	INFO                  // 普通信息
	WARN                  // 警告信息
	ERROR                 // 错误信息
	FATAL                 // 致命错误，会导致程序退出
)

var (
	currentLevel LogLevel    = INFO // 默认日志级别
	logger       *log.Logger        // 标准库日志实例
	once         sync.Once          // 确保初始化只执行一次
)

// initLogger 初始化日志器。
func initLogger() {
	once.Do(func() {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	})
}

// SetLogLevel 设置当前的日志级别。
func SetLogLevel(level LogLevel) {
	currentLevel = level
}

// Any 是一个辅助结构，用于记录键值对。
type Any struct {
	Key   string
	Value interface{}
}

// String 是一个辅助结构，用于记录字符串键值对。
type String struct {
	Key   string
	Value string
}

// Int 是一个辅助结构，用于记录整数键值对。
type Int struct {
	Key   string
	Value int
}

// Duration 是一个辅助结构，用于记录时间间隔键值对。
type Duration struct {
	Key   string
	Value time.Duration
}

// formatArgs 将可变参数转换为日志字符串。
func formatArgs(args []interface{}) string {
	if len(args) == 0 {
		return ""
	}

	var formatted string
	for _, arg := range args {
		switch v := arg.(type) {
		case Any:
			formatted += fmt.Sprintf(" %s=%v", v.Key, v.Value)
		case String:
			formatted += fmt.Sprintf(" %s=%s", v.Key, v.Value)
		case Int:
			formatted += fmt.Sprintf(" %s=%d", v.Key, v.Value)
		case Duration:
			formatted += fmt.Sprintf(" %s=%v", v.Key, v.Value)
		default:
			formatted += fmt.Sprintf(" %v", v)
		}
	}
	return formatted
}

// Debug 打印调试日志。
func Debug(msg string, args ...interface{}) {
	initLogger()
	if currentLevel <= DEBUG {
		logger.Printf("[DEBUG] %s%s\n", msg, formatArgs(args))
	}
}

// Info 打印信息日志。
func Info(msg string, args ...interface{}) {
	initLogger()
	if currentLevel <= INFO {
		logger.Printf("[INFO] %s%s\n", msg, formatArgs(args))
	}
}

// Warn 打印警告日志。
func Warn(msg string, args ...interface{}) {
	initLogger()
	if currentLevel <= WARN {
		logger.Printf("[WARN] %s%s\n", msg, formatArgs(args))
	}
}

// Error 打印错误日志。
func Error(msg string, args ...interface{}) {
	initLogger()
	if currentLevel <= ERROR {
		logger.Printf("[ERROR] %s%s\n", msg, formatArgs(args))
	}
}

// Fatal 打印致命错误日志并退出程序。
func Fatal(msg string, args ...interface{}) {
	initLogger()
	logger.Printf("[FATAL] %s%s\n", msg, formatArgs(args))
	os.Exit(1)
}
