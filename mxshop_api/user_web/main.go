package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mxshop_api/user_web/global"
	"mxshop_api/user_web/initialize"
	"mxshop_api/user_web/utils"
	myvalidator "mxshop_api/user_web/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化 logger
	initialize.InitLogger()

	// 2. 初始化配置文件
	initialize.InitConfig()

	// 3. 初始化 routers
	Router := initialize.Routers()

	// 4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 5. 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// ==========================================
	// 🚀 新增：Consul 注册逻辑
	// ==========================================

	// 6. 添加健康检查接口
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 7. 连接 Consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	consulClient, err := api.NewClient(cfg)
	if err != nil {
		panic("连接 Consul 失败: " + err.Error())
	}

	// 8. 获取本机真实 IP (Docker 内网 IP)
	localIP, err := utils.GetOutBoundIP()
	if err != nil {
		panic("获取本机 IP 失败: " + err.Error())
	}

	// 9. 生成服务注册对象
	serviceID := fmt.Sprintf("%s-%s", global.ServerConfig.Name, uuid.NewV4().String())
	port := global.ServerConfig.Port

	// 配置健康检查 (HTTP 方式)
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", localIP, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s",
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    global.ServerConfig.Name,
		Port:    port,
		Tags:    global.ServerConfig.Tags,
		Address: localIP,
		Check:   check,
	}

	// 10. 执行注册
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		panic("服务注册失败: " + err.Error())
	}
	zap.S().Infof("服务注册成功，ID: %s, 地址: %s:%d", serviceID, localIP, port)

	// ==========================================
	// 🚀 启动服务 & 优雅退出逻辑
	// ==========================================

	// 11. 在 goroutine 中启动 Gin 服务 (非阻塞)
	go func() {
		zap.S().Debugf("启动服务器，端口: %d", port)
		zap.S().Infof("【调试】当前完整配置: %+v", global.ServerConfig)
		if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	// 12. 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 主线程阻塞在这里，直到收到信号

	// 13. 收到信号，执行注销逻辑
	if err := consulClient.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("服务注销失败")
	} else {
		zap.S().Info("服务注销成功")
	}
}