package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
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

//TimeFormat 默认时间合并
const TimeFormat string = "2006-01-02 15:04:05"

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
