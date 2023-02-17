package gorm

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	GDb *gorm.DB
)

func Init() error {
	cfg, err := ini.Load("../../conf/ini/my.ini")
	if err != nil {
		log.Fatalf("加载配置文件失败。[err]=%s", err.Error())
	}
	section := cfg.Section("mysql")

	dsn := section.Key("dsn").String()
	maxIdle := section.Key("max_idle_conn").MustInt(5)
	maxOpen := section.Key("max_open_conn").MustInt(50)
	conn, err := connection(dsn, maxIdle, maxOpen)
	if err != nil {
		return fmt.Errorf("连接失败。[err]=%s", err.Error())
	}

	GDb = conn
	return nil
}

func connection(dsn string, idle, open int) (*gorm.DB, error) {
	connDB, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "blog_", // 表前缀

		},
	})
	if err != nil {
		return nil, err
	}
	db, err := connDB.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(open)

	return connDB, nil
}
