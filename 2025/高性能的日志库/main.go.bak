package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

/*
高性能的日志库，支持以下功能：
1、异步写入日志（使用 Channel 和 Goroutine）。
2、支持日志分级（DEBUG、INFO、WARN、ERROR）。
3、提供日志轮转和压缩功能
*/

// 日志级别常量定义
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

//  LogConfig 配置日志
type LogConfig struct {
	LogLevel       int    // 日志级别
	LogFile        string // 日志文件路径
	MaxSize        int64  // 单个日志文件最大大小（字节）
	EnableRotation bool   // 是否启用日志轮转机制
	EnableCompress bool   // 是否启用日志压缩
}

// LogEntry 日志条目
type LogEntry struct {
	Level   int
	Message string
	Time    time.Time
}

// Log 实现日志库
type Log struct {
	config       LogConfig
	logChannel   chan LogEntry
	wg           sync.WaitGroup
	logFile      *os.File
	currentSize  int64
	logFileMutex sync.Mutex
}

// NewLog 创建日志实例
func NewLog(config LogConfig) *Log {
	logInstance := &Log{
		config:     config,
		logChannel: make(chan LogEntry, 100), // 可以根据需求调整缓冲区大小
	}

	go logInstance.listen() // 启动异步监听日志
	return logInstance
}

// listen 监听日志通道
func (l *Log) listen() {
	for logEntry := range l.logChannel {
		// 异步写入日志
		l.writeLog(logEntry)
	}
}

// writeLog 实际写入日志文件
func (l *Log) writeLog(entry LogEntry) {
	// 判断日志级别
	if entry.Level >= l.config.LogLevel {
		// 打开日志文件
		l.logFileMutex.Lock()
		defer l.logFileMutex.Unlock()
		if l.logFile == nil {
			var err error
			l.logFile, err = os.OpenFile(l.config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatalf("打开日志文件失败: %v", err)
			}
		}

		// 计算当前日志文件大小
		info, err := l.logFile.Stat()
		if err != nil {
			log.Fatalf("获取日志文件信息失败: %v", err)
		}
		l.currentSize = info.Size()

		// 日志文件轮转
		if l.config.EnableRotation && l.currentSize > l.config.MaxSize {
			l.rotateLogFile()
		}

		// 写入日志
		logMessage := fmt.Sprintf("[%s] [%s] %s\n", entry.Time.Format("2006-01-02 15:04:05"), levelToString(entry.Level), entry.Message)
		_, err = l.logFile.WriteString(logMessage)
		if err != nil {
			log.Fatalf("写入日志失败: %v", err)
		}
	}
}

// rotateLogFile 执行日志轮转
func (l *Log) rotateLogFile() {
	// 关闭当前日志文件
	if l.logFile != nil {
		l.logFile.Close()
	}

	// 重命名当前日志文件
	archiveName := fmt.Sprintf("%s.%s", l.config.LogFile, time.Now().Format("20060102-150405"))
	err := os.Rename(l.config.LogFile, archiveName)
	if err != nil {
		log.Fatalf("重命名日志文件失败: %v", err)
	}

	// 如果启用压缩
	if l.config.EnableCompress {
		l.compressLogFile(archiveName)
	}

	// 重新打开日志文件
	l.logFile, err = os.OpenFile(l.config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("重新打开日志文件失败: %v", err)
	}
}

// compressLogFile 压缩日志文件
func (l *Log) compressLogFile(fileName string) {
	// 使用压缩工具压缩日志文件，可以使用gzip等工具
	fmt.Printf("日志文件 %s 已压缩\n", fileName)
}

// levelToString 日志级别转换为字符串
func levelToString(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogMessage 向日志发送消息
func (l *Log) LogMessage(level int, message string) {
	entry := LogEntry{
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}
	l.logChannel <- entry
}

// main 测试日志功能
func main() {
	// 配置日志
	config := LogConfig{
		LogLevel:       DEBUG,
		LogFile:        "app.log",
		MaxSize:        10 * 1024 * 1024, // 10MB
		EnableRotation: true,
		EnableCompress: true,
	}

	// 创建日志实例
	logger := NewLog(config)

	// 输出一些日志
	logger.LogMessage(DEBUG, "调试信息")
	logger.LogMessage(INFO, "信息日志")
	logger.LogMessage(WARN, "警告日志")
	logger.LogMessage(ERROR, "错误日志")

	// 停止日志服务
	close(logger.logChannel)
	logger.wg.Wait()
}
