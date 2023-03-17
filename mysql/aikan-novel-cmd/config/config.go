package config

type Cfg struct {
	MySQLCfg MySQLConfig
}

type MySQLConfig struct {
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
}
