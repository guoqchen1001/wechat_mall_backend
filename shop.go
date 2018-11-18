//shop.go 店铺
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

// Banner 首页banner结构
type Banner struct {
	gorm.Model
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
