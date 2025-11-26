package initialize

import "go.uber.org/zap"

// 通用代码，配置logger
func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
