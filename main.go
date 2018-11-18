// wechat_mall_backend project main.go
package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

//ResponseData 通用信息返回结构
type ResponseData struct {
	Code int         `json:"code"`    // 返回值
	Msg  string      `json:"message"` // 错误信息
	Data interface{} `json:"data"`    // 返回数据
}

func main() {

	mux := httprouter.New()
	mux.GET("/config/get-value", GetValue)                 // 获取配置
	mux.GET("/user/wxapp/login", WXAppLogin)               // 小程序登录
	mux.GET("/user/wxapp/register/complex", WXAppRegister) // 小程序注册
	mux.GET("/user/check-token", CheckToken)               // 校验token
	mux.GET("/score/send/rule", ScoreSendRule)             // 积分赠送规则
	mux.GET("/shop/banner/list", GetBannerList)            // 获取banner

	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}
	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	log.WithFields(logrus.Fields{
		"webserver": "error",
	}).Error(err.Error)

}
