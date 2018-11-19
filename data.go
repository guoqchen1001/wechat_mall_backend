// data.go 通用数据结构
package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

//ResponseData 通用信息返回结构
type ResponseData struct {
	Code int         `json:"code"`    // 返回值
	Msg  string      `json:"message"` // 错误信息
	Data interface{} `json:"data"`    // 返回数据
}

//WxConfig 微信配置信息，由参数表获取
type WxConfig struct {
	AppID     string
	AppSecret string
	GrantType string
	Code      string
}

//WxSessionResponse 微信获取session返回信息
type WxSessionResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrorCode  int    `json:"errcode"`     // 错误码
	ErrorMsg   string `json:"errmsg"`      // 错误信息
}

//TimeFormat 默认时间合并
const TimeFormat string = "2006-01-02 15:04:05"

//Init 获取微信小程序基本配置
func (wxConfig *WxConfig) Init() error {

	var config Config

	// 获取小程序id
	config = Config{}
	db.Where("No = ?", "appId").First(&config)

	if config != (Config{}) {
		wxConfig.AppID = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	// 获取小程序密钥
	config = Config{}
	db.Where("No = ?", "appSecret").First(&config)

	if config != (Config{}) {
		wxConfig.AppSecret = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	wxConfig.GrantType = "authorization_code"

	return nil
}

var db *gorm.DB
var err error

// init 数据库连接
func init() {

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=wechat dbname=wechat_mall password=123 sslmode=disable")

	if err != nil {
		log.WithFields(logrus.Fields{
			"db": "connect",
		}).Panic(err)
	}

}

// init 数据库建表
func init() {
	// 创建基础配置表
	err = db.AutoMigrate(&Config{}, &User{}, &Banner{}, &Category{}).Error
	if err != nil {
		log.WithFields(logrus.Fields{
			"db": "initTable",
		}).Panic(err)
	}
}

//init 基础数据，后续需转化接口配置
func init() {

	// 写入基础数据-店铺名称
	config := new(Config)
	config.No = "mallName"
	config.Val = "小卖铺"
	config.Name = "店铺名称"

	entry := log.WithFields(logrus.Fields{"config": "init"})

	err := db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据-小程序appid
	config = new(Config)
	config.No = "appId"
	config.Val = "wxb05d528592b74609"
	config.Name = "小程序appid"

	err = db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据库-小程序密钥
	config = new(Config)
	config.No = "appSecret"
	config.Val = "d89d3ee840cd9dd89015086962229f52"
	config.Name = "小程序密钥"

	err = db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据-最低充值金额
	config = new(Config)
	config.No = "recharge_amount_min"
	config.Val = "1"
	config.Name = "充值最少金额"

	err = db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

}

// init 写入banner数据，后续虚改为上传
func init() {
	// 写入banner数据，需要接口上传
	banner := new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_1.jpg"
	banner.Order = 1
	banner.Status = "0"
	banner.Title = "1"
	banner.StatusStr = "显示"
	db.FirstOrCreate(&banner, banner)

	banner = new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_2.jpg"
	banner.Order = 2
	banner.Status = "0"
	banner.Title = "2"
	banner.StatusStr = "显示"
	db.FirstOrCreate(&banner, banner)

	banner = new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_3.jpg"
	banner.Order = 3
	banner.Status = "0"
	banner.Title = "3"
	banner.StatusStr = "显示"
	db.FirstOrCreate(&banner, banner)

}

// init 写入类别，后续需提供接口
func init() {

	entry := log.WithFields(logrus.Fields{
		"category": "init",
	})

	// 写入现切水果
	category := new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "01"
	category.Name = "现切水果"
	category.Order = 1
	err := db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入招牌推荐
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "02"
	category.Name = "招牌推荐"
	category.Order = 2
	err = db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入好吃推荐
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "03"
	category.Name = "好吃推荐"
	category.Order = 3
	err = db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入休闲干果
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "04"
	category.Name = "休闲干果"
	category.Order = 4
	err = db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

}
