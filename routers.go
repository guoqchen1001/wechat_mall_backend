package main

import (
	"net/http"
	"wechat_mall_backend/controllers"

	"github.com/julienschmidt/httprouter"
)

var mux http.Handler

// CeateRouter 创建路由
func CeateRouter() error {

	mux := httprouter.New()
	mux.GET("/config/get-value", controllers.GetValue)                 // 获取配置
	mux.GET("/user/wxapp/login", controllers.WXAppLogin)               // 小程序登录
	mux.GET("/user/wxapp/register/complex", controllers.WXAppRegister) // 小程序注册
	mux.GET("/user/check-token", controllers.CheckToken)               // 校验token
	mux.GET("/score/send/rule", controllers.ScoreSendRule)             // 积分赠送规则
	mux.GET("/banner/list", controllers.GetBannerList)                 // 获取banner
	mux.GET("/shop/goods/category/all", controllers.GetCategoryList)   // 获取类别
	mux.ServeFiles("/static/*filepath", http.Dir("static"))            // 静态文件

	return nil

}
