package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var Conf *AppConfig

type AppConfig struct {
	CfgFile    string
	Env        string
	HttpPort   string
	LogFile    string
	LogConsole bool
	LogLevel   string

	Raw      *ini.File
	DBConfig []DBConfig
}

type DBConfig struct {
	Name        string
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
}

// Load 加载配置文件
func Load(file string) (*AppConfig, error) {
	Conf = &AppConfig{
		CfgFile:  file,
		Raw:      ini.Empty(),
		Env:      Dev,
		HttpPort: "8080",
	}

	// 判断配置文件是否存在
	if _, err := os.Stat(Conf.CfgFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件[%s]不存在", Conf.CfgFile)
	}

	// 初始化配置
	conf, err := ini.Load(Conf.CfgFile)
	if err != nil {
		return nil, fmt.Errorf("初始化配置[%s]失败", Conf.CfgFile)
	}

	Conf.Raw = conf
	Conf.LoadAppCfg()
	Conf.LoadDBCfg()

	return Conf, nil
}

func (cfg *AppConfig) LoadAppCfg() {
	section := cfg.Raw.Section("app")
	if env := section.Key("env").String(); env != "" {
		cfg.Env = env
	}
	if httpPort := section.Key("http_port").String(); httpPort != "" {
		cfg.HttpPort = httpPort
	}
	if logFile := section.Key("log_file").String(); logFile != "" {
		cfg.LogFile = logFile
	}
	if logConsole := section.Key("log_console").String(); logConsole == "true" {
		cfg.LogConsole = true
	}
	if logLevel := section.Key("log_level").String(); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}

func (cfg *AppConfig) LoadDBCfg() {
	section := cfg.Raw.Section("db")
	cfg.DBConfig = []DBConfig{
		{
			Name:        "default",
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(5),
			MaxOpenConn: section.Key("max_open_conn").MustInt(10),
		},
	}

	childSections := cfg.Raw.Section("db").ChildSections()
	for _, section := range childSections {
		cfg.DBConfig = append(cfg.DBConfig, DBConfig{
			Name:        strings.TrimLeft(section.Name(), "db."),
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(5),
			MaxOpenConn: section.Key("max_open_conn").MustInt(10),
		})
	}
}

func (cfg *AppConfig) IsDevEnv() bool {
	return cfg.Env == "dev"
}
