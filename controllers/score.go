package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wechat_mall_backend/models"

	"github.com/julienschmidt/httprouter"
)

//ScoreSendRule 积分赠送规则
func ScoreSendRule(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {

			var responseData ResponseData
			responseData.Code = 0
			responseData.Msg = err.(error).Error()
			output, _ := json.Marshal(responseData)
			fmt.Fprint(w, string(output))
		}

	}()

	r.ParseForm()

	if len(r.Form["code"]) <= 0 {
		panic(fmt.Errorf("未找到参数%s", "code"))
	}

	var responseData ResponseData
	code := r.Form["code"][0]
	// 好评送
	if code == "goodReputation" {

		var scoreRuleList []models.ScoreRule

		scoreRule := models.ScoreRule{
			Code:    "goodReputation",
			CodeStr: "好评送",
			Confine: 0.00,
			Score:   3,
		}
		scoreRuleList = append(scoreRuleList, scoreRule)

		responseData.Code = 0
		responseData.Data = scoreRuleList

		output, err := json.Marshal(responseData)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, string(output))

	} else {
		panic(fmt.Errorf("未定义的类型%s", code))
	}

}
