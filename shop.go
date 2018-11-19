//shop.go 店铺
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

//OrgBanner 不带gorm.model的banner
type OrgBanner struct {
	BusinessID int    `json:"businessId"` // 商户ID
	LinkURL    string `json:"linkUrl"`    // 链接地址
	Order      int    `json:"paixu"`      // 排序
	PicURL     string `json:"picUrl"`     // 图片地址
	Name       string `json:"remark"`     // 名称
	Status     string `json:"status"`     // 状态
	StatusStr  string `json:"statusStr"`  // 状态名称
	Title      string `json:"title"`      // 标题
	Type       string `json:"type"`       // 类型
	UserID     int    `json:"userId"`     // 用户id
}

//Banner 首页banner
type Banner struct {
	gorm.Model
	OrgBanner
}

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

//MarshalJSON 自定义Banner的json序列化方法
func (banner Banner) MarshalJSON() ([]byte, error) {

	type TmpBanner struct {
		ID         uint   `json:"id"`
		DateAdd    string `json:"dataAdd"`
		DateUpdate string `json:"dateUpdate"`
		OrgBanner
	}

	tmp := TmpBanner{
		ID:         banner.ID,
		DateAdd:    banner.CreatedAt.Format(TimeFormat),
		DateUpdate: banner.UpdatedAt.Format(TimeFormat),
		OrgBanner:  banner.OrgBanner,
	}

	return json.Marshal(tmp)

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

//GetBannerList 返回banner列表
func GetBannerList(w http.ResponseWriter, r *http.Request, o httprouter.Params) {
	var banners []Banner
	var responseData ResponseData

	defer func() {
		if err := recover(); err != nil {
			var responseData ResponseData
			responseData.Code = 10
			responseData.Msg = err.(error).Error()

			output, _ := json.Marshal(responseData)
			fmt.Fprint(w, string(output))
		}
	}()

	err := db.Find(&banners).Error
	if err != nil {
		panic(err)
	}

	responseData.Code = 0
	responseData.Data = banners

	output, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))

}

// GetCategoryList 获取类别列表
func GetCategoryList(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var responsData ResponseData
			responsData.Code = 10
			responsData.Msg = err.(error).Error()
			output, _ := json.Marshal(responsData)
			fmt.Fprint(w, string(output))
		}
	}()

	var categoies []Category
	err := db.Find(&categoies).Error
	if err != nil {
		panic(err)
	}

	var responseData ResponseData
	responseData.Code = 0
	responseData.Data = categoies
	output, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))

}
