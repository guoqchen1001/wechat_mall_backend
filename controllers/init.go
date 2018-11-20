package controllers

//ResponseData 通用信息返回结构
type ResponseData struct {
	Code int         `json:"code"`    // 返回值
	Msg  string      `json:"message"` // 错误信息
	Data interface{} `json:"data"`    // 返回数据
}
