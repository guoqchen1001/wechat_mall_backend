package models

//ScoreRule 积分规则结构
type ScoreRule struct {
	Code    string  // 积分代码
	CodeStr string  `json:"codeStr"` // 积分字符说明
	Confine float32 // 适用范围
	Score   int     // 分值3
}
