// store.go
package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// 配置文件信息，对应数据库中参数配置表t_sys_parn
type Config struct {
	No   string `gorm:"primary_key" json:"no"`
	Val  string `json:"value"`
	Name string `json:"name"`
}

// 微信配置信息，由参数表获取
type WxConfig struct {
	AppId     string
	AppSecret string
	GrantType string
	Code      string
}

// 用户
type User struct {
	UserId    string `gorm:"primary_key"` // 用户id
	UserNo    string // 用户编码
	UserName  string // 用户名
	WechatId  string `json:"openId"` // 微信id
	Phone     string // 手机号
	Country   string // 国家
	NickName  string `json:"nickName"`  //  昵称
	AvatarUrl string `json:"avatarurl"` // 图像地址
	Gender    int    // 性别
	Province  string // 省份
	City      string // 城市
	Language  string // 显示语言

}

// 微信获取session返回信息
type WxSessionResponse struct {
	OpenId     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionId    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrorCode  int    `json:"errcode"`     // 错误码
	ErrorMsg   string `json:"errmsg"`      // 错误信息
}

func (wx_config *WxConfig) Init() error {

	var config Config

	// 获取小程序id
	config = Config{}
	db.Where("No = ?", "appId").First(&config)

	if config != (Config{}) {
		wx_config.AppId = config.Val
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
		panic(err)
	}
	// 创建基础配置表
	db.AutoMigrate(&Config{}, &User{})

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

}
