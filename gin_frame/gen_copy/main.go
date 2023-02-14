package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"toolbox/gin_frame/gen_copy/config"
	"toolbox/gin_frame/gen_copy/models/mysql"
	"toolbox/gin_frame/gen_copy/router"
	"toolbox/gin_frame/gen_copy/zlog"
)

func main() {
	// 加载启动配置文件
	var configFile string
	flag.StringVar(&configFile, "conf", "app.ini", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		panic(fmt.Sprintf("加载配置文件失败。[file]=%s;[ERR]=%s", configFile, err))
	}

	// 初始化日志
	zlog.Init(cfg)
	logger := zlog.GetLogger()
	defer logger.Sync()

	// 初始化数据库
	err = mysql.Init(cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化MySQL数据库失败。[ERR]=%s", err.Error()))
	}

	// 启动服务
	zlog.Info("服务启动中...")
	err = startServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("启动服务失败。[ERR]=%s", err.Error()))
	}

}

func startServer(cfg *config.AppConfig) error {
	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: getEngine(cfg),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctxFunc context.CancelFunc) {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT)
		for {
			select {
			case <-signalCh:
				ctxFunc()
				return
			}
		}
	}(cancel)
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			panic(fmt.Sprintf("停止GIN服务失败。[ERR]=%s", err.Error()))
		}
	}()
	log.Println("GIN服务启动成功...")

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		//log.Println("GIN服务已正常关闭")
		zlog.Debug("GIN服务已正常关闭")
		return nil
	}
	return err
}

func getEngine(cfg *config.AppConfig) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.IsDevEnv() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "服务器内部错误，请稍后再试",
		})
	}))

	// 初始化路由
	log.Println("加载路由...")
	router.RegisterRouters(engine)

	return engine
}
