// wechat_mall_backend project main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// web返回配置信息
type ConfigData struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data Config `json:"data"`
}

func get_value(w http.ResponseWriter, r *http.Request, o httprouter.Params) {
	r.ParseForm()

	var config Config
	var config_data ConfigData

	if len(r.Form["key"]) > 0 {

		db.Where("No = ?", r.Form["key"]).First(&config)

		if config != (Config{}) {
			config_data.Data = config
			config_data.Code = 0
			config_data.Msg = ""
		} else {
			config_data.Code = 10
			config_data.Msg = fmt.Sprintf("尚未设置[%s]", r.Form["key"])
		}

	} else {
		config_data.Code = 20
		config_data.Msg = fmt.Sprintf("未找到请求参数[%s]", "key")
	}

	output, err := json.Marshal(config_data)

	fmt.Println(string(output))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, string(output))
	}

}

func main() {

	mux := httprouter.New()
	mux.GET("/config/get-value", get_value)

	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}
	server.ListenAndServeTLS("cert.pem", "key.pem")

}
