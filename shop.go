//shop.go 店铺
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

type Banner struct {
	gorm.Model
	BusinessId int    `json:"businessId"` // 商户ID
	Id         int    `json:"id"`         // id
	LinkUrl    string `json:"linkUrl"`    // 链接地址
	Order      int    `json:"paixu"`      // 排序
	PicUrl     string `json:"picUrl"`     // 图片地址
	Name       string `json:"remark"`     // 名称
	Status     string `json:"status"`     // 状态
	StatusStr  string `json:"statusStr"`  // 状态名称
	Title      string `json:"title"`      // 标题
	Type       string `json:"type"`       // 类型
	UserId     int    `json:"userId"`     // 用户id
}

// 返回banner列表
func GetBannerList(w http.ResponseWriter, r *http.Request, o httprouter.Params) {
	var banners []Banner
	var response_data ResponseData

	defer func() {
		if err := recover(); err != nil {
			var response_data ResponseData
			response_data.Code = 10
			response_data.Msg = err.(error).Error()

			output, _ := json.Marshal(response_data)
			fmt.Fprint(w, string(output))
		}
	}()

	err := db.Find(&banners).Error
	if err != nil {
		panic(err)
	}

	response_data.Code = 0
	response_data.Data = banners

	output, err := json.Marshal(response_data)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))

}
