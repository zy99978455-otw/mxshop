package initialize

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8" //新版替换
	"go.uber.org/zap"
	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)

	cfg := elasticsearch.Config{
		Addresses: []string{host},
	}
	var err error
	// 创建官方客户端
	global.EsClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("创建 ES 客户端失败: %s", err.Error()))
	}

	// 验证连接是否可用 (可选，但推荐)
	info, err := global.EsClient.Info()
	if err != nil {
		panic(fmt.Sprintf("连接 ES 失败: %s", err.Error()))
	}
	defer info.Body.Close()
	zap.S().Infof("ES 连接成功，版本信息: %s", info.String())

	// 2. 检查索引是否存在
	indexName := model.EsGoods{}.GetIndexName()
	
	// 官方库的 Exists 方法返回 Response，需要根据 StatusCode 判断
	res, err := global.EsClient.Indices.Exists([]string{indexName})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// StatusCode 200 表示存在，404 表示不存在
	if res.StatusCode == 404 {
		zap.S().Infof("索引 %s 不存在，正在创建...", indexName)
		
		// 3. 创建 Mapping 和 Index
		// 官方库需要传入 io.Reader，所以使用 strings.NewReader
		createRes, err := global.EsClient.Indices.Create(
			indexName,
			global.EsClient.Indices.Create.WithBody(strings.NewReader(model.EsGoods{}.GetMapping())),
		)
		if err != nil {
			panic(err)
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			panic(fmt.Sprintf("创建索引失败: %s", createRes.String()))
		}
		zap.S().Infof("索引 %s 创建成功", indexName)
	} else {
		zap.S().Infof("索引 %s 已存在，跳过创建", indexName)
	}
}
