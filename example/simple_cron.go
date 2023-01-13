package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron"
)

func main() {
	client := cron.New()
	_ = client.AddFunc("*/2 * * * * *", func() {
		fmt.Println("[timer]=", time.Now().String())
	})
	client.Start()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
