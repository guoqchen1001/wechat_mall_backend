// wechat_mall_backend project main.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 通用信息返回结构
type ResponseData struct {
	Code int         `json:"code"`    // 返回值
	Msg  string      `json:"message"` // 错误信息
	Data interface{} `json:"data"`    // 返回数据
}

func main() {

	mux := httprouter.New()
	mux.GET("/config/get-value", get_value)
	mux.GET("/user/wxapp/login", wxapp_login)
	mux.GET("/user/wxapp/register/complex", wxapp_register)

	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}
	server.ListenAndServeTLS("cert.pem", "key.pem")

}
