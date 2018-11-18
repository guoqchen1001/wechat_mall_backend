// data.go 通用数据结构
package main

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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

//Init 获取微信小程序基本配置
func (wx_config *WxConfig) Init() error {

	var config Config

	// 获取小程序id
	config = Config{}
	db.Where("No = ?", "appId").First(&config)

	if config != (Config{}) {
		wx_config.AppID = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	// 获取小程序密钥
	config = Config{}
	db.Where("No = ?", "appSecret").First(&config)

	if config != (Config{}) {
		wx_config.AppSecret = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	wx_config.GrantType = "authorization_code"

	return nil
}

var db *gorm.DB
var err error

func init() {

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=wechat dbname=wechat_mall password=123 sslmode=disable")

	if err != nil {
		log.WithFields(logrus.Fields{
			"db": "connect",
		}).Panic(err)
	}
	// 创建基础配置表
	err = db.AutoMigrate(&Config{}, &User{}, &Banner{}).Error
	if err != nil {
		log.WithFields(logrus.Fields{
			"db": "init",
		}).Panic(err)
	}

	// 基础数据后续需转化为sql语句执行

	// 写入基础数据-店铺名称
	config := new(Config)
	config.No = "mallName"
	config.Val = "小卖铺"
	config.Name = "店铺名称"

	db.FirstOrCreate(&config, config)

	// 写入基础数据-小程序appid
	config = new(Config)
	config.No = "appId"
	config.Val = "wxb05d528592b74609"
	config.Name = "小程序appid"

	db.FirstOrCreate(&config, config)

	// 写入基础数据库-小程序密钥
	config = new(Config)
	config.No = "appSecret"
	config.Val = "d89d3ee840cd9dd89015086962229f52"
	config.Name = "小程序密钥"

	db.FirstOrCreate(&config, config)

	// 写入基础数据-最低充值金额
	config = new(Config)
	config.No = "recharge_amount_min"
	config.Val = "1"
	config.Name = "充值最少金额"

	db.FirstOrCreate(&config, config)

}
