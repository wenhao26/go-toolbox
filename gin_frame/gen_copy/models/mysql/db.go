package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"toolbox/gin_frame/gen_copy/config"
)

var (
	db          *gorm.DB
	connCluster map[string]*gorm.DB
)

func Init(cfg *config.AppConfig) error {
	connCluster := make(map[string]*gorm.DB)
	for _, v := range cfg.DBConfig {
		conn, err := Connection(v.Dsn, v.MaxIdleConn, v.MaxOpenConn)
		if err != nil {
			return fmt.Errorf("连接失败。[ERR]=%s", err.Error())
		}

		connCluster[v.Name] = conn
		if v.Name == "default" {
			db = conn
		}
	}
	return nil
}

func Connection(dsn string, idle, open int) (*gorm.DB, error) {
	connDB, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}))
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

func DB(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := connCluster[dbName[0]]; ok {
			return conn
		}
	}
	return db
}
