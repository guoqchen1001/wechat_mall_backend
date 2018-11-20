package models

import (
	"wechat_mall_backend/logs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //pg驱动
	"github.com/sirupsen/logrus"
)

// Db 数据库对象
var Db *gorm.DB
var err error

// init 数据库连接
func init() {

	Db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=wechat dbname=wechat_mall password=123 sslmode=disable")

	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"db": "connect",
		}).Panic(err)
	}

}

// init 数据库建表
func init() {
	// 创建基础配置表
	err := Db.AutoMigrate(Config{}, User{}, Banner{}, Category{}).Error
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"db": "initTable",
		}).Panic(err)
	}
}

//init 基础数据，后续需转化接口配置
func init() {

	// 写入基础数据-店铺名称
	config := new(Config)
	config.No = "mallName"
	config.Val = "小卖铺"
	config.Name = "店铺名称"

	entry := logs.Log.WithFields(logrus.Fields{"config": "init"})

	err := Db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据-小程序appid
	config = new(Config)
	config.No = "appId"
	config.Val = "wxb05d528592b74609"
	config.Name = "小程序appid"

	err = Db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据库-小程序密钥
	config = new(Config)
	config.No = "appSecret"
	config.Val = "d89d3ee840cd9dd89015086962229f52"
	config.Name = "小程序密钥"

	err = Db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入基础数据-最低充值金额
	config = new(Config)
	config.No = "recharge_amount_min"
	config.Val = "1"
	config.Name = "充值最少金额"

	err = Db.FirstOrCreate(&config, config).Error
	if err != nil {
		entry.Error(err.Error)
	}

}

// init 写入banner数据，后续虚改为上传
func init() {
	// 写入banner数据，需要接口上传
	banner := new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_1.jpg"
	banner.Order = 1
	banner.Status = "0"
	banner.Title = "1"
	banner.StatusStr = "显示"
	Db.FirstOrCreate(&banner, banner)

	banner = new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_2.jpg"
	banner.Order = 2
	banner.Status = "0"
	banner.Title = "2"
	banner.StatusStr = "显示"
	Db.FirstOrCreate(&banner, banner)

	banner = new(Banner)
	banner.PicURL = "https://localhost:8081/static/banner_3.jpg"
	banner.Order = 3
	banner.Status = "0"
	banner.Title = "3"
	banner.StatusStr = "显示"
	Db.FirstOrCreate(&banner, banner)

}

// init 写入类别，后续需提供接口
func init() {

	entry := logs.Log.WithFields(logrus.Fields{
		"category": "init",
	})

	// 写入现切水果
	category := new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "01"
	category.Name = "现切水果"
	category.Order = 1
	err := Db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入招牌推荐
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "02"
	category.Name = "招牌推荐"
	category.Order = 2
	err = Db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入好吃推荐
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "03"
	category.Name = "好吃推荐"
	category.Order = 3
	err = Db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

	// 写入休闲干果
	category = new(Category)
	category.Level = 1
	category.IsUse = true
	category.No = "04"
	category.Name = "休闲干果"
	category.Order = 4
	err = Db.FirstOrCreate(&category, category).Error
	if err != nil {
		entry.Error(err.Error)
	}

}
