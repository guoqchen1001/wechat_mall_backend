// user.go 用户
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/xlstudio/wxbizdatacrypt"
)

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

// 微信登录
func WXAppLogin(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

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

	err = db.Where(&User{WechatId: wx_session_response.OpenId}).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if user == (User{}) {
		panic(UserNotFoundError)
	}

	token, err := CreateToken(user.UserId)
	if err != nil {
		panic(err)
	}

	type LoginData struct {
		Uid   string `json:"uid"`
		Token string `json:"token"`
	}

	var login_data LoginData
	login_data.Uid = user.UserId
	login_data.Token = token

	response_data.Code = 0
	response_data.Data = login_data

	output, err := json.Marshal(response_data)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

// 微信用户注册
func WXAppRegister(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

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

	type RegisterData struct {
		UserId string `json:"user_id"`
	}

	var register_data RegisterData
	register_data.UserId = user.UserId
	response_data.Code = 0
	response_data.Data = register_data

	output, _ := json.Marshal(response_data)
	fmt.Fprint(w, string(output))

}

// 检查token是否合法
func CheckToken(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var response_data ResponseData
			response_data.Code = 10
			response_data.Msg = err.(error).Error()
			output, _ := json.Marshal(response_data)
			fmt.Fprint(w, output)
		}
	}()

	r.ParseForm()
	if len(r.Form["token"]) <= 0 {
		panic(fmt.Errorf("未找到有效的%s值", "token"))
	}

	token := r.Form["token"][0]
	_, err := ParseToken(token)
	if err != nil {
		panic(err)
	}

	var respons_data ResponseData
	respons_data.Code = 0
	output, err := json.Marshal(respons_data)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, output)
	return

}
