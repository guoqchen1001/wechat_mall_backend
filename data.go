// data.go 通用数据结构
package main

import (
	"errors"
	"wechat_mall_backend/models"
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

//Init 获取微信小程序基本配置
func (wxConfig *WxConfig) Init() error {

	type Config models.Config
	var config Config

	// 获取小程序id
	config = Config{}
	models.Db.Where("No = ?", "appId").First(&config)

	if config != (Config{}) {
		wxConfig.AppID = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	// 获取小程序密钥
	config = Config{}
	models.Db.Where("No = ?", "appSecret").First(&config)

	if config != (Config{}) {
		wxConfig.AppSecret = config.Val
	} else {
		return errors.New("未找到有效的小程序appid，请检查系统设置")
	}

	wxConfig.GrantType = "authorization_code"

	return nil
}
