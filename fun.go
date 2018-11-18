// fun.go 通用函数
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
