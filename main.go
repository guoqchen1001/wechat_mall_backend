// wechat_mall_backend project main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/xlstudio/wxbizdatacrypt"
)

// web返回配置信息
type ResponseData struct {
	Code int    `json:"code"`    // 返回值
	Msg  string `json:"message"` // 错误信息
	Data Config `json:"data"`    // 返回数据
}

// 获取系统参数值
func get_value(w http.ResponseWriter, r *http.Request, o httprouter.Params) {
	r.ParseForm()

	var config Config
	var response_data ResponseData

	if len(r.Form["key"]) > 0 {

		db.Where("No = ?", r.Form["key"]).First(&config)

		if config != (Config{}) {
			response_data.Data = config
			response_data.Code = 0
			response_data.Msg = ""
		} else {
			response_data.Code = 10
			response_data.Msg = fmt.Sprintf("尚未设置[%s]", r.Form["key"])
		}

	} else {
		response_data.Code = 20
		response_data.Msg = fmt.Sprintf("未找到请求参数[%s]", "key")
	}

	output, err := json.Marshal(response_data)

	fmt.Println(string(output))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, string(output))
	}

}

// 微信登录
func wxapp_login(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	var response_data ResponseData
	var wx_session_response WxSessionResponse
	var (
		UserNotFoundError = errors.New("用户不存在")
		CodeNotFoundError = errors.New("请传入有效的code值")
	)

	defer func() {
		if err := recover(); err != nil {

			switch err {
			case UserNotFoundError:
				response_data.Code = 10000

			case CodeNotFoundError:
				response_data.Code = 20

			default:
				response_data.Code = 10
			}

			response_data.Msg = err.(error).Error()
			output, _ := json.Marshal(response_data)
			fmt.Fprint(w, string(output))
		}
	}()

	// 解析参数
	r.ParseForm()

	if len(r.Form["code"]) <= 0 {
		panic(CodeNotFoundError)
	}

	// 获取微信session和openid
	code := r.Form["code"][0]
	wx_session_response, err = Code2Session(code)
	if err != nil {
		panic(err)
	}

	// 用户不存在需要注册
	user := User{}
	db.Where("WechatId = ?", wx_session_response.OpenId).First(&user)
	if user == (User{}) {
		panic(UserNotFoundError)
	}

}

// 微信用户注册
func wxapp_register(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var response_data ResponseData // 最终返回值
			response_data.Code = 10
			response_data.Msg = err.(error).Error()
			output, _ := json.Marshal(response_data)
			fmt.Fprint(w, string(output))
		}
	}()

	// 解析参数
	r.ParseForm()

	var response_data ResponseData            // 最终返回值
	var wx_session_response WxSessionResponse // 通过小程序code获取的session值
	var wx_config WxConfig                    // 微信设置

	params := []string{
		"code", "encryptedData", "iv",
	}

	for _, param := range params {
		if len(r.Form[param]) <= 0 {
			panic(errors.New(fmt.Sprintf("请传入有效的%s值", param)))
		}
	}

	code := r.Form["code"][0]
	encryptedData := r.Form["encryptedData"][0]
	iv := r.Form["iv"][0]

	// 获取session
	wx_session_response, err = Code2Session(code)
	if err != nil {
		panic(err)
	}

	// 获取微信设置
	err = wx_config.Init()
	if err != nil {
		panic(err)
	}

	pc := wxbizdatacrypt.WxBizDataCrypt{
		AppID:      wx_config.AppId,
		SessionKey: wx_session_response.SessionKey,
	}

	result, err := pc.Decrypt(encryptedData, iv, true)
	if err != nil {
		panic(err)
	}

	// 生成uuid作为userid
	userid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	user := new(User)
	user.UserId = userid.String()

	result_byte := []byte(result.(string))
	err = json.Unmarshal(result_byte, &user)
	if err != nil {
		panic(err)
	}

	err = db.Create(user).Error
	if err != nil {
		panic(err)
	}

	if db.NewRecord(user) {
		panic(errors.New("保存用户到数据库失败"))
	}

	response_data.Code = 0
	response_data.Data = Config{
		No:   "userid",
		Name: "用户id",
		Val:  user.UserId,
	}
	output, _ := json.Marshal(response_data)
	fmt.Fprint(w, string(output))

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
