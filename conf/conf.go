package conf

import (
	"gopkg.in/ini.v1"
)

func GetINI() *ini.File {
	cfg, err := ini.Load("../conf/ini/my.ini")
	if err != nil {
		panic(err)
	}
	return cfg
}
