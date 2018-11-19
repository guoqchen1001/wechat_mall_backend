// scr0e.go 积分
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//ScoreRule 积分规则结构
type ScoreRule struct {
	Code    string  // 积分代码
	CodeStr string  `json:"codeStr"` // 积分字符说明
	Confine float32 // 适用范围
	Score   int     // 分值3
}

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
		panic(err)
	}

	var responseData ResponseData
	code := r.Form["code"][0]
	// 好评送
	if code == "goodReputation" {

		var scoreRuleList []ScoreRule

		scoreRule := ScoreRule{
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
