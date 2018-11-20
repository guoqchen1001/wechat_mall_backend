package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"wechat_mall_backend/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/xlstudio/wxbizdatacrypt"
)

//WxConfig 微信配置信息，由参数表获取
type WxConfig struct {
	AppID     string
	AppSecret string
	GrantType string
	Code      string
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

//WxSessionResponse 微信获取session返回信息
type WxSessionResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrorCode  int    `json:"errcode"`     // 错误码
	ErrorMsg   string `json:"errmsg"`      // 错误信息
}

const wxCode2sessionURL = "https://api.weixin.qq.com/sns/jscode2session"
const tokenSecretKey = "guoqchen"

//Code2Session 通过小程序的code获取登录session
// 入参 code string 小程序登录时微信返回值
// 回参 WxSessionResponse  error
func Code2Session(code string) (WxSessionResponse, error) {

	client := http.Client{}

	var wxConfig WxConfig
	var wxSessionResponse WxSessionResponse

	req, err := http.NewRequest("GET", wxCode2sessionURL, nil)
	if err != nil {
		return wxSessionResponse, err
	}

	// 微信小程序登录获取的code
	wxConfig.Code = code

	// 获取微信小程序配置, appid,appsecret,grant_type
	err = wxConfig.Init()
	if err != nil {
		return wxSessionResponse, err
	}

	params := req.URL.Query()
	params.Add("appid", wxConfig.AppID)
	params.Add("secret", wxConfig.AppSecret)
	params.Add("js_code", wxConfig.Code)
	params.Add("grant_type", wxConfig.GrantType)

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return wxSessionResponse, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wxSessionResponse, err
	}

	err = json.Unmarshal(body, &wxSessionResponse)
	if err != nil {
		return wxSessionResponse, err
	}

	return wxSessionResponse, nil

}

// CreateToken 创建登录Token
// 入参 userID string
// 回参 userID string, err error
func CreateToken(userID string) (string, error) {
	signKey := []byte(tokenSecretKey)

	type CustomClaim struct {
		UserID string
		jwt.StandardClaims
	}

	exp := time.Now().Add(24 * time.Hour)
	claims := CustomClaim{
		userID,
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// ParseToken 解析Token
// 入参 tokenStr string
// 回参 userid string   err error
func ParseToken(tokenStr string) (string, error) {
	type CustomClaim struct {
		UserID string
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", err

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
	wxSessionResponse, err := Code2Session(code)
	if err != nil {
		panic(err)
	}

	// 用户不存在需要注册
	user := models.User{}

	err = models.Db.Where(&models.User{WechatID: wxSessionResponse.OpenID}).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if user == (models.User{}) {
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
	wxSessionResponse, err := Code2Session(code)
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
	user := new(models.User)
	user.UserID = userid.String()

	resultByte := []byte(result.(string))
	err = json.Unmarshal(resultByte, &user)
	if err != nil {
		panic(err)
	}

	err = models.Db.Create(user).Error
	if err != nil {
		panic(err)
	}

	if models.Db.NewRecord(user) {
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
