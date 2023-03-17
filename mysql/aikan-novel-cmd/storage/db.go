package storage

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"toolbox/mysql/aikan-novel-cmd/global"
)

func InitDB() {
	config := global.Cfg.MySQLCfg
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: config.Dsn}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	sqlDb.SetMaxOpenConns(config.MaxOpenConn)

	global.DB = db
}
