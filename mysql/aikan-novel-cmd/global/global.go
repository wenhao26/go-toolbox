package global

import (
	"gorm.io/gorm"

	"toolbox/mysql/aikan-novel-cmd/config"
)

var (
	Cfg config.Cfg
	DB  *gorm.DB
)
