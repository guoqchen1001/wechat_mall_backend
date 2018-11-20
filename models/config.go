package models

//Config 配置文件信息
type Config struct {
	No   string `gorm:"primary_key" json:"no"`
	Val  string `json:"value"`
	Name string `json:"remark"`
}
