package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type S2sLog struct {
	Id int `json:"id"`
}

var (
	delLimit  int
	cronSpec  int
	runNumber int

	db *gorm.DB
)

func init() {
	config, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	BaseSection := config.Section("base")
	delLimit = BaseSection.Key("del_limit").MustInt(500)
	cronSpec = BaseSection.Key("cron_spec").MustInt(30)
	runNumber = BaseSection.Key("run_number").MustInt(5)

	MysqlSection := config.Section("mysql")
	dsn := MysqlSection.Key("dsn").String()
	maxIdle := MysqlSection.Key("max_idle_conn").MustInt(5)
	maxOpen := MysqlSection.Key("max_open_conn").MustInt(50)

	DB, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent), // 最低级，无论如何都不输出日志了
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)

	db = DB
}

func (S2sLog) TableName() string {
	return "iw_s2s_log"
}

func delLogs() {
	var s2sLog S2sLog

	s2sLogs := []S2sLog{}
	db.Debug().Table(s2sLog.TableName()).Select("id").Limit(delLimit).Find(&s2sLogs)
	if len(s2sLogs) == 0 {
		fmt.Println("无处理数据，退出程序")
		os.Exit(0)
	}
	db.Table(s2sLog.TableName()).Delete(&s2sLogs)
	log.Println("删除日志数量：", len(s2sLogs))
}

func main() {
	p, _ := ants.NewPool(10, ants.WithPreAlloc(false))
	defer p.Release()

	c := cron.New()
	_ = c.AddFunc(fmt.Sprintf("*/%d * * * * *", cronSpec), func() {
		for i := 0; i < runNumber; i++ {
			_ = p.Submit(func() {
				delLogs()
			})
		}
	})
	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
