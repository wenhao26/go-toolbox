package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 日志库
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)

	/*logrus.Trace("trade msg")
	logrus.Debug("debug msg")
	logrus.Info("info msg")
	logrus.Warn("warn msg")
	logrus.Error("error msg")
	logrus.Fatal("fatal msg")
	logrus.Panic("panic msg")*/

	// 重定向输出
	w1 := &bytes.Buffer{}
	w2 := os.Stdout
	w3, err := os.OpenFile("data/log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	logrus.SetOutput(io.MultiWriter(w1, w2, w3))
	logrus.Info("info msg")

}
