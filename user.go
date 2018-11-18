// user.go 用户
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/xlstudio/wxbizdatacrypt"
)

//User 用户
type User struct {
	UserID    string `gorm:"primary_key"` // 用户id
	UserNo    string // 用户编码
	UserName  string // 用户名
	WechatID  string `json:"openId"` // 微信id
	Phone     string // 手机号
	Country   string // 国家
	NickName  string `json:"nickName"`  //  昵称
	AvatarURL string `json:"avatarurl"` // 图像地址
	Gender    int    // 性别
	Province  string // 省份
	City      string // 城市
	Language  string // 显示语言

}

//WXAppLogin 微信登录
func WXAppLogin(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	var responseData ResponseData
	var wxSessionResponse WxSessionResponse
	var (
		UserNotFoundError = errors.New("用户不存在")
		CodeNotFoundError = errors.New("请传入有效的code值")
	)

	defer func() {
		if err := recover(); err != nil {

			switch err {
			case UserNotFoundError:
				responseData.Code = 10000

			case CodeNotFoundError:
				responseData.Code = 20

			default:
				responseData.Code = 10
			}

			responseData.Msg = err.(error).Error()
			output, _ := json.Marshal(responseData)
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
	wxSessionResponse, err = Code2Session(code)
	if err != nil {
		panic(err)
	}

	// 用户不存在需要注册
	user := User{}

	err = db.Where(&User{WechatID: wxSessionResponse.OpenID}).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if user == (User{}) {
		panic(UserNotFoundError)
	}

	token, err := CreateToken(user.UserID)
	if err != nil {
		panic(err)
	}

	type LoginData struct {
		UID   string `json:"uid"`
		Token string `json:"token"`
	}

	var loginData LoginData
	loginData.UID = user.UserID
	loginData.Token = token

	responseData.Code = 0
	responseData.Data = loginData

	output, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

//WXAppRegister 微信用户注册
func WXAppRegister(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var responseData ResponseData // 最终返回值
			responseData.Code = 10
			responseData.Msg = err.(error).Error()
			output, _ := json.Marshal(responseData)
			fmt.Fprint(w, string(output))
		}
	}()

	// 解析参数
	r.ParseForm()

	var responseData ResponseData           // 最终返回值
	var wxSessionResponse WxSessionResponse // 通过小程序code获取的session值
	var wxConfig WxConfig                   // 微信设置

	params := []string{
		"code", "encryptedData", "iv",
	}

	for _, param := range params {
		if len(r.Form[param]) <= 0 {
			panic(fmt.Errorf("请传入有效的%s值", param))
		}
	}

	code := r.Form["code"][0]
	encryptedData := r.Form["encryptedData"][0]
	iv := r.Form["iv"][0]

	// 获取session
	wxSessionResponse, err = Code2Session(code)
	if err != nil {
		panic(err)
	}

	// 获取微信设置
	err = wxConfig.Init()
	if err != nil {
		panic(err)
	}

	pc := wxbizdatacrypt.WxBizDataCrypt{
		AppID:      wxConfig.AppID,
		SessionKey: wxSessionResponse.SessionKey,
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
	user.UserID = userid.String()

	resultByte := []byte(result.(string))
	err = json.Unmarshal(resultByte, &user)
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
		UserID string `json:"user_id"`
	}

	var registerData RegisterData
	registerData.UserID = user.UserID
	responseData.Code = 0
	responseData.Data = registerData

	output, _ := json.Marshal(responseData)
	fmt.Fprint(w, string(output))

}

//CheckToken 检查token是否合法
func CheckToken(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var responseData ResponseData
			responseData.Code = 10
			responseData.Msg = err.(error).Error()
			output, _ := json.Marshal(responseData)
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

	var responsData ResponseData
	responsData.Code = 0
	output, err := json.Marshal(responsData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, output)
	return

}
