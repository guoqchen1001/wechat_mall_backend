//shop.go 店铺
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

// 不带gorm.model的banner
type OrgBanner struct {
	BusinessId int    `json:"businessId"` // 商户ID
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

// 首页banner
type Banner struct {
	gorm.Model
	OrgBanner
}

// 自定义Banner的json序列化方法
func (this Banner) MarshalJSON() ([]byte, error) {

	type TmpBanner struct {
		ID         uint   `json:"id"`
		DateAdd    string `json:"dataAdd"`
		DateUpdate string `json:"dateUpdate"`
		OrgBanner
	}

	tmp := TmpBanner{
		ID:         this.ID,
		DateAdd:    this.CreatedAt.Format(TimeFormat),
		DateUpdate: this.UpdatedAt.Format(TimeFormat),
		OrgBanner:  this.OrgBanner,
	}

	return json.Marshal(tmp)

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
