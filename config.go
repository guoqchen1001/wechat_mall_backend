package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

// 配置文件信息
type Config struct {
	No   string `gorm:"primary_key" json:"no"`
	Val  string `json:"value"`
	Name string `json:"name"`
}

// 获取系统参数值
func get_value(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	var config Config
	var response_data ResponseData

	defer func() {
		if err := recover(); err != nil {
			response_data.Code = 10
			response_data.Msg = err.(error).Error()
			output, _ := json.Marshal(response_data)
			fmt.Fprint(w, string(output))
		}

	}()

	r.ParseForm()

	if len(r.Form["key"]) > 0 {

		err := db.Where("No = ?", r.Form["key"]).First(&config).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}

		if config != (Config{}) {
			response_data.Code = 0
			response_data.Data = config

		} else {
			panic(errors.New(fmt.Sprintf("尚未设置[%s]，请检查系统设置", r.Form["key"])))
		}

	} else {
		panic(errors.New(fmt.Sprintf("未找到请求参数%s", "key")))
	}

	output, err := json.Marshal(response_data)

	if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(w, string(output))
	}

}
