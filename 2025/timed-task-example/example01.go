package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var c = cron.New(cron.WithSeconds())

func addTask(spec string, taskFunc func()) (cron.EntryID, error) {
	entryID, err := c.AddFunc(spec, taskFunc)
	if err != nil {
		return 0, err
	}

	c.Start()

	return entryID, nil
}

func stopTask(entryID cron.EntryID, name string) {
	c.Remove(entryID)

	fmt.Printf(">>>> %s stopped\n", name)
}

func main() {
	entryID1, err := addTask("*/3 * * * * *", func() {
		fmt.Println("task01 running at", time.Now())
	})
	if err != nil {
		panic(err)
	}

	entryID2, err := addTask("*/5 * * * * *", func() {
		fmt.Println("task02 running at", time.Now())
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(entryID1, entryID2)

	time.Sleep(10 * time.Second)

	stopTask(entryID2, "task02")

	// 每隔10秒重置定时器
	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			timer.Reset(10 * time.Second)
		}
	}
}
