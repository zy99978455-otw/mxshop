package global

import( 
	"mxshop_api/user_web/config"
	"mxshop_api/user_web/proto"
	ut "github.com/go-playground/universal-translator"
)
var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	// 🚨【新增】存放本地 config-pro.yaml 的内容 (如: Nacos地址, 8848端口)
    NacosConfig *config.NacosConfig = &config.NacosConfig{}
	UserSrvClient proto.UserClient	
)