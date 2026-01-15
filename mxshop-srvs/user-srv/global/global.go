package global

import (
	"gorm.io/gorm"
	"mxshop-srvs/user-srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
