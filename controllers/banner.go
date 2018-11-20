package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wechat_mall_backend/models"

	"github.com/julienschmidt/httprouter"
)

//GetBannerList 返回banner列表
func GetBannerList(w http.ResponseWriter, r *http.Request, o httprouter.Params) {
	var banners []models.Banner
	var responseData ResponseData

	defer func() {
		if err := recover(); err != nil {
			var responseData ResponseData
			responseData.Code = 10
			responseData.Msg = err.(error).Error()

			output, _ := json.Marshal(responseData)
			fmt.Fprint(w, string(output))
		}
	}()

	err := models.Db.Find(&banners).Error
	if err != nil {
		panic(err)
	}

	responseData.Code = 0
	responseData.Data = banners

	output, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))

}
