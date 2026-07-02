package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(level, output string) error {
	var zl zapcore.Level

	if err := zl.UnmarshalText([]byte(level)); err != nil {
		zl = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zl),
		Encoding:         "json",
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}

	// 控制台输出使用 console 格式更易读，但保持 JSON 便于日志采集
	if output == "stdout" || output == "stderr" {
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}

func Sync() {
	_ = Log.Sync()
}
