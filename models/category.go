package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

//OrgCategory 类别
type OrgCategory struct {
	Icon   string `json:"icon"`   // 图标
	IsUse  bool   `json:"isUse"`  // 是否使用
	No     string `json:"key"`    // 类别编码
	Level  int    `json:"level"`  // 等级
	Name   string `json:"name"`   // 名称
	Order  int    `json:"paixu"`  // 排序
	PID    int    `json:"pid"`    // pid
	Type   string `json:"type"`   // 类型
	UserID int    `json:"userId"` // 用户id
}

// Category 类别，gorm类型
type Category struct {
	gorm.Model
	OrgCategory
}

// MarshalJSON 自定义json输出
func (category Category) MarshalJSON() ([]byte, error) {
	type TmpCategory struct {
		ID      uint   `json:"id"`
		DateAdd string `json:"dateAdd"`
		OrgCategory
	}

	tmp := TmpCategory{
		ID:          category.ID,
		DateAdd:     category.CreatedAt.Format(TimeFormat),
		OrgCategory: category.OrgCategory,
	}
	return json.Marshal(tmp)
}
