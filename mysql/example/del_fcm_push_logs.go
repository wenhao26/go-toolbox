package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type FcmPushLog struct {
	Id int `json:"id"`
	//PushId     int    `json:"push_id"`
	//UserId     int    `json:"user_id"`
	//PushData   string `json:"push_data"`
	//PushResult string `json:"push_result"`
	//Status     int    `json:"status"`
	//IsTest     int    `json:"is_test"`
	//CreateTime int    `json:"create_time"`
}

var (
	delLimit int
	cronSpec int

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

	MysqlSection := config.Section("mysql")
	dsn := MysqlSection.Key("dsn").String()
	maxIdle := MysqlSection.Key("max_idle_conn").MustInt(5)
	maxOpen := MysqlSection.Key("max_open_conn").MustInt(50)

	DB, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

func (FcmPushLog) TableName() string {
	return "iw_fcm_push_log"
}

func delLogs() {
	var pushLog FcmPushLog

	pushLogs := []FcmPushLog{}
	db.Table(pushLog.TableName()).Select("id").Limit(delLimit).Find(&pushLogs)
	if len(pushLogs) == 0 {
		//os.Exit(0)
	}
	db.Table(pushLog.TableName()).Delete(&pushLogs)
	log.Println("删除日志数量：", len(pushLogs))
}

func main() {
	c := cron.New()
	_ = c.AddFunc(fmt.Sprintf("*/%d * * * * *", cronSpec), func() {
		go delLogs()
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
