package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wechat_mall_backend/models"

	"github.com/julienschmidt/httprouter"
)

// GetCategoryList 获取类别列表
func GetCategoryList(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			var responsData ResponseData
			responsData.Code = 10
			responsData.Msg = err.(error).Error()
			output, _ := json.Marshal(responsData)
			fmt.Fprint(w, string(output))
		}
	}()

	var categoies []models.Category
	err := models.Db.Find(&categoies).Error
	if err != nil {
		panic(err)
	}

	var responseData ResponseData
	responseData.Code = 0
	responseData.Data = categoies
	output, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))

}
