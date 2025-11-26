package router

import (
	"mxshop_api/user_web/api"
	"mxshop_api/user_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList) //单个组件执行middlewares
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}