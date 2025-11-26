package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
	"mxshop_srv/goods_srv/global"
)

//类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为不为null
//https://zhuanlan.zhihu.com/p/73997266
//这些类型我们使用int32还是int

// 商品目录
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null;comment:分类名称" json:"name"`
	Level            int32       `gorm:"type:int;not null;default:1;comment:分类级别" json:"level"`
	IsTab            bool        `gorm:"default:false;not null;comment:是否在Tab标签展示" json:"is_tab"`
	// 核心自关联字段
	ParentCategoryID int32       `gorm:"comment:父分类的ID" json:"parent"`
	ParentCategory   *Category   `gorm:"comment:父分类对象" json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID;comment:子分类列表" json:"sub_category"`

}

// 商品品牌
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;comment:品牌名称"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:品牌logo"`
}
// 中间表，处理多对多关系
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:轮播图URI"`
	Url   string `gorm:"type:varchar(200);not null;comment:点击轮播图跳转地址"`
	Index int32  `gorm:"type:int;default:1;not null;comment:轮播图显示顺序(数值越小越靠前)"`
}

// 商品信息
type Goods struct {
	// 基础信息和关联
	BaseModel
	CategoryID int32 `gorm:"type:int;not null;comment:商品分类ID"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null;comment:品牌ID"`
	Brands     Brands
	// 状态和属性
	OnSale   bool `gorm:"default:false;not null;comment:是否在售"`
	ShipFree bool `gorm:"default:false;not null;comment:是否包邮"`
	IsNew    bool `gorm:"default:false;not null;comment:是否新品"`
	IsHot    bool `gorm:"default:false;not null;comment:是否热门/爆款"`
	// 核心商品信息
	Name       string `gorm:"type:varchar(50);not null;comment:商品名称"`
	GoodsSn    string `gorm:"type:varchar(50);not null;comment:商品序列号"`
	GoodsBrief string `gorm:"type:varchar(100);not null;comment:商品简介"`
	// 统计数据
	ClickNum int32 `gorm:"type:int;default:0;not null;comment:商品被点击/浏览的数量"`
	SoldNum  int32 `gorm:"type:int;default:0;not null;comment:商品累计销量"`
	FavNum   int32 `gorm:"type:int;default:0;not null;comment:商品收藏量"`
	// 价格信息
	MarketPrice float32 `gorm:"not null;comment:市场指导价(原价)"`
	ShopPrice   float32 `gorm:"not null;comment:实际销售价"`
	// 图片信息
	Images          GormList `gorm:"type:varchar(1000);not null;comment:商品轮播图URI"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:商品详情描述图片列表URI"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:商品主图URI"`
}

// -----------------------------------------------------------------------------
// GORM Hooks: 同步数据到 Elasticsearch (v8 官方客户端写法)
// 

// AfterCreate 在 MySQL 创建成功后，将商品数据写入 ES
func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	// 1. 构造 ES 模型数据
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	// 2. 序列化为 JSON
	data, err := json.Marshal(esModel)
	if err != nil {
		return err
	}

	// 3. 构建 Index 请求 (对应旧版 client.Index().BodyJson(...))
	req := esapi.IndexRequest{
		Index:      esModel.GetIndexName(),
		DocumentID: strconv.Itoa(int(g.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true", // 强制刷新，确保数据立即可查
	}

	// 4. 执行请求
	res, err := req.Do(context.Background(), global.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 5. 检查业务错误
	if res.IsError() {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		return fmt.Errorf("ES Index Error: %s", buf.String())
	}
	return nil
}

// AfterUpdate 在 MySQL 更新成功后，更新 ES 中的文档
func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	// 构造更新体，使用 "doc" 包装（局部更新）
	updateBody := map[string]interface{}{
		"doc": esModel,
	}
	data, err := json.Marshal(updateBody)
	if err != nil {
		return err
	}

	// 构建 Update 请求
	req := esapi.UpdateRequest{
		Index:      esModel.GetIndexName(),
		DocumentID: strconv.Itoa(int(g.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), global.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		return fmt.Errorf("ES Update Error: %s", buf.String())
	}
	return nil
}



// AfterDelete 在 MySQL 删除成功后，从 ES 中删除文档
func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	// 构建 Delete 请求
	req := esapi.DeleteRequest{
		Index:      EsGoods{}.GetIndexName(),
		DocumentID: strconv.Itoa(int(g.ID)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), global.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 忽略 404 Not Found 错误 (可能 ES 里本来就没有)
	if res.IsError() && !strings.Contains(res.Status(), "Not Found") {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		return fmt.Errorf("ES Delete Error: %s", buf.String())
	}
	return nil
}