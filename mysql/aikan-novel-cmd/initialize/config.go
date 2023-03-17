package initialize

import (
	"gopkg.in/ini.v1"

	"toolbox/mysql/aikan-novel-cmd/config"
	"toolbox/mysql/aikan-novel-cmd/global"
)

func InitConfig() {
	f, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	mysqlCfg := mysqlConfig(f)

	global.Cfg = config.Cfg{MySQLCfg: mysqlCfg}
}

func mysqlConfig(f *ini.File) config.MySQLConfig {
	section := f.Section("mysql")
	dsn := section.Key("dsn").String()
	maxIdle := section.Key("max_idle_conn").MustInt(5)
	maxOpen := section.Key("max_open_conn").MustInt(50)

	return config.MySQLConfig{
		Dsn:         dsn,
		MaxIdleConn: maxIdle,
		MaxOpenConn: maxOpen,
	}
}
