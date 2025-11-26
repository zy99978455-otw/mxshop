package initialize

import (
	"fmt"
	"strings"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"go.uber.org/zap"

	"mxshop_srv/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

func InitConfig(){
	//从配置文件中读取出对应的配置
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"

	var configFileName string

	// 判断环境变量来选择配置文件
	
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	} else {
		configFileName = fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	}
	
	fmt.Println("Using config file:", configFileName) //输入文件路径，调试用

	// 创建 viper 实例
	v := viper.New()
	v.SetConfigFile(configFileName)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//将配置文件内容加载到全局变量,读取Nacos配置
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	// 【新增代码】: 打印原始值（用 %q 查看是否有 \r），并进行清洗
    zap.S().Infof("清洗前 Nacos 配置: Host=%q, Namespace=%q", global.NacosConfig.Host, global.NacosConfig.Namespace)

    // 强制清洗 Host, Namespace, DataId, Group 中的空格和换行符
    global.NacosConfig.Host = strings.TrimSpace(global.NacosConfig.Host)
    global.NacosConfig.Namespace = strings.TrimSpace(global.NacosConfig.Namespace)
    global.NacosConfig.DataId = strings.TrimSpace(global.NacosConfig.DataId)
    global.NacosConfig.Group = strings.TrimSpace(global.NacosConfig.Group)

    zap.S().Infof("清洗后 Nacos 配置: Host=%q, Namespace=%q", global.NacosConfig.Host, global.NacosConfig.Namespace)
	zap.S().Infof("获取本地 Nacos 配置: %+v", global.NacosConfig)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port: global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig {
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	// 获取远程配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}
// 【修正 2】：远程配置通常是业务配置，这时才加载到 global.ServerConfig
    // 注意：通常 Nacos 存的是 YAML 格式，建议用 yaml.Unmarshal 而不是 json
    // 如果你确定 Nacos 上存的是 JSON，请保留 json.Unmarshal
    // err = json.Unmarshal([]byte(content), &global.ServerConfig)
	err = yaml.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil{
		zap.S().Fatalf("解析nacos配置失败： %s", err.Error())
	}
}