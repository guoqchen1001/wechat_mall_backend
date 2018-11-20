// wechat_mall_backend project main.go
package main

import (
	"net/http"
	"os"
	"path/filepath"
	"wechat_mall_backend/controllers"

	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

func main() {

	mux := httprouter.New()
	mux.GET("/config/get-value", GetValue)                 // 获取配置
	mux.GET("/user/wxapp/login", WXAppLogin)               // 小程序登录
	mux.GET("/user/wxapp/register/complex", WXAppRegister) // 小程序注册
	mux.GET("/user/check-token", CheckToken)               // 校验token
	mux.GET("/score/send/rule", controllers.ScoreSendRule) // 积分赠送规则
	mux.GET("/banner/list", GetBannerList)                 // 获取banner
	mux.GET("/shop/goods/category/all", GetCategoryList)   // 获取类别

	mux.ServeFiles("/static/*filepath", http.Dir("static"))

	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.WithFields(logrus.Fields{
			"webserver": "error",
		}).Error(err)
	}

	certPem := filepath.Join(dir, "pem/cert.pem")
	keyPem := filepath.Join(dir, "pem/key.pem")

	err = server.ListenAndServeTLS(certPem, keyPem)

	log.WithFields(logrus.Fields{
		"webserver": "error",
	}).Error(err)

}
