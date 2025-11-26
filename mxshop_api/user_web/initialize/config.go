package initialize

import (
	"fmt"
	"mxshop_api/user_web/global"
	"gopkg.in/yaml.v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// 🚨 必须引入这三个 Nacos SDK 包
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant" // <--- 之前缺的就是这行！
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// ==============================
	// Step 1: 读取本地配置 (Bootstrap)
	// ==============================
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	// 修正点：去掉 "mxshop_api/user_web/" 前缀
    // 因为程序运行的当前目录里直接就有 config-pro.yaml
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	// 文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 🚨【关键修正 1】这里读取的是 Nacos 的连接信息，要存到 NacosConfig 里！
	// 不能存到 ServerConfig，否则 ServerConfig.Port 会变成 8848
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("已加载本地 Nacos 配置: %+v", global.NacosConfig)

	// ==============================
	// Step 2: 连接 Nacos 读取远程配置
	// ==============================
	// 1. 创建 Nacos 客户端配置
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port, //8848
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		// RotateTime:          "1h",
		// MaxAge:              3,
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	// 读取远程配置 (DataID: user-web.yaml)
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic("连接 Nacos 获取配置失败: " + err.Error())
	}

	// 4. 解析远程配置到 ServerConfig
	// 注意：这里要把 content (JSON/YAML字符串) 解析到 global.ServerConfig
	// 如果 Nacos 上是 YAML 格式，这里其实建议用 yaml.Unmarshal，但如果为了兼容之前的习惯：
	// 假设 Nacos 内容也是 JSON 格式，或者用 viper 解析字符串：
	
    // 推荐方式：直接解析 Nacos 返回的字符串
    // 假设你在 Nacos 存的是 YAML，需要 import "gopkg.in/yaml.v3"
    // err = yaml.Unmarshal([]byte(content), &global.ServerConfig)
    
    // 如果你在 Nacos 存的是 JSON:
	err = yaml.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic("解析 Nacos 配置失败: " + err.Error())
	}

	zap.S().Infof("🚀 成功加载远程配置，启动端口: %d", global.ServerConfig.Port)
    
    // (可选) 这里可以加 ListenConfig 动态监听逻辑...
}