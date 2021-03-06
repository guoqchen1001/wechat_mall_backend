package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wechat_mall_backend/models"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

//GetValue 获取系统参数值
func GetValue(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	var config models.Config
	var responseData ResponseData

	defer func() {
		if err := recover(); err != nil {
			responseData.Code = 10
			responseData.Msg = err.(error).Error()
			output, _ := json.Marshal(responseData)
			fmt.Fprint(w, string(output))
		}

	}()

	r.ParseForm()

	if len(r.Form["key"]) > 0 {

		err := models.Db.Where("No = ?", r.Form["key"]).First(&config).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}

		if config != (models.Config{}) {
			responseData.Code = 0
			responseData.Data = config

		} else {
			panic(fmt.Errorf("尚未设置[%s]，请检查系统设置", r.Form["key"]))
		}

	} else {
		panic(fmt.Errorf("未找到请求参数%s", "key"))
	}

	output, err := json.Marshal(responseData)

	if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(w, string(output))
	}

}
