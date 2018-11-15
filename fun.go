package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const wx_code2session_url = "https://api.weixin.qq.com/sns/jscode2session"
const token_secret_key = "guoqchen"

// 通过小程序的code获取登录session
func Code2Session(code string) (WxSessionResponse, error) {

	client := http.Client{}

	var wx_config WxConfig
	var wx_session_response WxSessionResponse

	req, err := http.NewRequest("GET", wx_code2session_url, nil)
	if err != nil {
		return wx_session_response, err
	}

	// 微信小程序登录获取的code
	wx_config.Code = code

	// 获取微信小程序配置, appid,appsecret,grant_type
	err = wx_config.Init()
	if err != nil {
		return wx_session_response, err
	}

	params := req.URL.Query()
	params.Add("appid", wx_config.AppId)
	params.Add("secret", wx_config.AppSecret)
	params.Add("js_code", wx_config.Code)
	params.Add("grant_type", wx_config.GrantType)

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return wx_session_response, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wx_session_response, err
	}

	err = json.Unmarshal(body, &wx_session_response)
	if err != nil {
		return wx_session_response, err
	}

	return wx_session_response, nil

}

func CreateToken(user_id string) (string, error) {
	sign_key := []byte(token_secret_key)

	type CustomClaim struct {
		UserId string
		jwt.StandardClaims
	}

	exp := time.Now().Add(24 * time.Hour)
	claims := CustomClaim{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(sign_key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ParseToken(token_str string) (string, error) {
	type CustomClaim struct {
		UserId string
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(token_str, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return "", err
	}

}
