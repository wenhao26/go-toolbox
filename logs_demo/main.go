package main

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogData struct {
	Platform     int     `json:"platform"`
	UserID       int     `json:"user_id"`
	Lang         string  `json:"lang"`
	AppVersion   string  `json:"app_version"`
	IP           string  `json:"ip"`
	Time         int64   `json:"time"`
	BookID       int     `json:"book_id"`
	ChapterID    int     `json:"chapter_id"`
	AdPlatform   string  `json:"ad_platform"`
	AdSpaceType  string  `json:"ad_space_type"`
	AdUnit       string  `json:"ad_unit"`
	AdFormat     string  `json:"ad_format"`
	AdSource     string  `json:"ad_source"`
	AdEarnings   float64 `json:"ad_earnings"`
	CurrencyCode string  `json:"currency_code"`
}

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()

	timer := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-timer.C:
			// output JSON
			jsons, err := json.Marshal(genLogData())
			if err != nil {
				panic(err.Error())
			}
			sugarLogger.Info(string(jsons))
			fmt.Println("output log...")

			timer.Reset(time.Second * 5)
		}
	}
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/test.log",
		MaxSize:    1,
		MaxAge:     30,
		MaxBackups: 5,
		LocalTime:  false,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func genLogData() LogData {
	return LogData{
		Platform:     1,
		UserID:       123456,
		Lang:         "en",
		AppVersion:   "3.8.8",
		IP:           "113.103.20.139",
		Time:         time.Now().Unix(),
		BookID:       1234,
		ChapterID:    1,
		AdPlatform:   "admob",
		AdSpaceType:  "task_ads",
		AdUnit:       "fa7750afe64f94ab",
		AdFormat:     "rewarded",
		AdSource:     "mintegral",
		AdEarnings:   0.00023852,
		CurrencyCode: "USD",
	}
}
