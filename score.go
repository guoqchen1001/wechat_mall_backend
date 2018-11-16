// scr0e.go 积分
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 积分规则结构
type ScoreRule struct {
	Code    string
	CodeStr string `json:"codeStr"`
	Confine float32
	Score   int
}

// 积分赠送规则
func ScoreSendRule(w http.ResponseWriter, r *http.Request, o httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {

			var response_date ResponseData
			response_date.Code = 0
			response_date.Msg = err.(error).Error()
			output, _ := json.Marshal(response_date)
			fmt.Fprint(w, string(output))
		}

	}()

	r.ParseForm()

	if len(r.Form["code"]) <= 0 {
		panic(err)
	}

	var response_data ResponseData
	code := r.Form["code"][0]
	// 好评送
	if code == "goodReputation" {

		var score_rule_list []ScoreRule

		score_rule := ScoreRule{
			Code:    "goodReputation",
			CodeStr: "好评送",
			Confine: 0.00,
			Score:   3,
		}
		score_rule_list = append(score_rule_list, score_rule)

		response_data.Code = 0
		response_data.Data = score_rule_list

		output, err := json.Marshal(response_data)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, string(output))

	} else {
		panic(fmt.Errorf("未定义的类型%s", code))
	}

}
