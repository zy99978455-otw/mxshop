package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// 自定义string切片，用于存放图片URI
type GormList []string

// Go 切片被序列化成 JSON 字符串并存入数据库的 VARCHAR 字段
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g) //转换为JSON格式的字节数组
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
// 数据库的 JSON 字符串被反序列化成 Go 切片，方便在业务代码中使用。
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

// 基础模板，解决重复
type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"` //为什么使用int32， bigint
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"` //json:"-" 是一种常用的方式，用于将内部管理字段与对外暴露的 API 数据隔离开来。
}
