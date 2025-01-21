package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// 实现一个高性能的日志系统，支持异步写入和日志分级。

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

// Logger 日志结构体
type Logger struct {
	file    *os.File
	logChan chan string
	wg      sync.WaitGroup
}

// NewLogger 创建日志实例
func NewLogger(filename string) *Logger {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	logger := &Logger{
		file:    file,
		logChan: make(chan string, 100),
	}

	logger.wg.Add(1)
	go logger.processLogs()
	return logger
}

// processLogs 处理日志
func (l *Logger) processLogs() {
	defer l.wg.Done()
	for logEntry := range l.logChan {
		_, _ = l.file.WriteString(logEntry)
	}
}

// Log 生成日志并写入管道
func (l *Logger) Log(level LogLevel, message string) {
	levelStr := ""
	switch level {
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}
	l.logChan <- fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), levelStr, message)
}

// Close 关闭
func (l *Logger) Close() {
	close(l.logChan)
	l.wg.Wait()
	l.file.Close()
}

func main() {
	logger := NewLogger("app.log")
	defer logger.Close()

	logger.Log(INFO, "INFO - 您可以对系统的整体状况和功能获得有意义的洞见")
	logger.Log(WARN, "WARN - 您可以获得事件时间表，以便更快地进行故障排查")
	logger.Log(ERROR, "ERROR - 您可以识别安全漏洞并将安全风险降到最低")

	fmt.Println("OK")
}
