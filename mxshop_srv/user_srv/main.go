package main

import (
	"flag"
	"fmt"
	"mxshop_srv/user_srv/global"
	"mxshop_srv/user_srv/handler"
	"mxshop_srv/user_srv/initialize"
	"mxshop_srv/user_srv/proto"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop_srv/user_srv/utils"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	ConsulAddr := flag.String("consul", "mxshop_consul:8500", "Consul服务地址")
	flag.Parse()

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info(global.ServerConfig)

	//创建 GRPC 服务
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 注册 gRPC 健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// ==========================================
	// 核心修改 1: 获取 Docker 容器的真实内网 IP
	// ==========================================
	// 此时得到的 ServiceAddr 类似 "172.20.0.5"
	// 如果 utils 包里还没有 GetOutBoundIP，请参考上一条回复补上
	ServiceAddr, err := utils.GetOutBoundIP()
	if err != nil {
		panic("获取本机IP失败: " + err.Error())
	}
	zap.S().Infof("检测到本机 IP: %s", ServiceAddr)

	// Consul服务注册
	cfg := api.DefaultConfig()
	cfg.Address = *ConsulAddr

	client, err := api.NewClient(cfg)
	if err != nil {
		panic("failed to create client: " + err.Error())
	}

	// 服务健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ServiceAddr, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "1m", //延长一点方便调试
	}

	serviceID := fmt.Sprintf("%s-%s", global.ServerConfig.Name, ServiceAddr) // 生成唯一的ID

	// 服务注册对象
	registration := &api.AgentServiceRegistration{
		ID:      serviceID, //ID 必须唯一
		Name:    global.ServerConfig.Name,
		Port:    *Port,
		Tags:    []string{"imooc", "bobby", "user", "srv"},
		Address: ServiceAddr, //使用真实IP
		Check:   check,
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		panic("failed to register service to consul: " + err.Error())
	}

	zap.S().Infof("服务已注册到Consul: %s:%d", ServiceAddr, *Port)

	// ==========================================
	// 核心修改 2: 优雅退出与信号监听
	// ==========================================
	// 启动 gRPC 服务 (放入协程，不阻塞主线程)
	go func() {
			if err := server.Serve(lis); err != nil {
				panic("failed to start grpc server: " + err.Error())
			}
		}()

	// 监听退出信号 (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在这里，直到收到信号

	// 收到退出信号，执行注销逻辑
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	} else {
		zap.S().Info("注销成功")
	}
}
