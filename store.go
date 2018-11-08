// store.go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Config struct {
	No   string `gorm:"primary_key" json:"no"`
	Val  string `json:"value"`
	Name string `json:"name"`
}

func (Config) TableName() string {
	return "t_sys_parm"
}

var db *gorm.DB
var err error

func init() {

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=wechat dbname=wechat_mall password=123 sslmode=disable")

	if err != nil {
		panic(err)
	}
	// 创建基础配置表
	db.AutoMigrate(&Config{})

	// 写入基础数据
	config := new(Config)
	config.No = "mallName"
	config.Val = "小卖铺"
	config.Name = "店铺名称"

	db.FirstOrCreate(&config, config)

}
