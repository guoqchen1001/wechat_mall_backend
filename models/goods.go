package models

import (
	"time"
)

// Goods 商品结构
type Goods struct {
	Barcode            string    `json:"barcode"`            //条码
	CategoryID         int       `json:"categoryid"`         //类别ID
	Name               string    `json:"name"`               //名称
	Characteristic     string    `json:"characteristic"`     //特色描述
	ShopIO             int       `json:"shopId"`             //店铺ID
	Order              int       `json:"paixu"`              // 排序
	IsGroup            bool      `json:"pintuan"`            //是否拼团
	IsBargain          bool      `json:"kanjian"`            //是否砍价
	LogisticID         int       `json:"logisticId"`         //物流模板
	Status             int       `json:"status"`             // 状态
	StartSellDate      time.Time `json:"dateStart"`          // 开始销售时间
	EndSellDate        time.Time `json:"dateEnd"`            // 停止销售时间
	Score              int       `json:"gotScore"`           // 获得积分
	ScoreType          int       `json:"gotScoreType"`       // 积分类型
	MinScore           int       `json:"minScore"`           // 最低积分
	RecommendStatus    int       `json:"recommendStatus"`    // 推荐状态
	RecommendStatusStr string    `json:"recommendStatusStr"` // 推荐状态类型
	Weight             float32   `json:"weight"`             // 重量

}

// GoodsNumber 商品指标数
type GoodsNumber struct {
	GoodsID              int  `json:"goodsID"`              // 商品ID
	NumberFav            uint `json:"numberfav"`            // 喜爱数
	NumberGoodReputation uint `json:"numberGoodReputation"` // 好评数
	NumberOrders         uint `json:"numberOrders"`         // 预定数
	NumberSells          uint `json:"numberSells"`          // 好评数
	ViewsNum             uint `json:"view"`                 // 浏览量
}

// GoodsPrice 商品价格
type GoodsPrice struct {
	GoodsID       int     `json:"goodsID"`       // 商品ID
	GoodsBarcode  string  `json:"barcode"`       // 商品价格
	OriginalPrice float32 `json:"originalPrice"` // 出初始价格
	MinPrice      float32 `jsonL:"minPrice"`     // 最低价格
	GroupPrice    float32 `json:"pintuanPrice"`  // 拼团价格
	BargainPrice  float32 `json:"kanjiaPrice"`   // 砍价价格
}

// GoodsPicture 商品图片
type GoodsPicture struct {
	GoodsID int    `json:"goodsId"` //商品ID
	PicURL  string `json:"pic"`     // 图片地址
	Order   string `json:"paixu"`   // 排序
}

//GoodsStock 商品库存
type GoodsStock struct {
	GoodsID int     `json:"goodsId"` // 商品ID
	Stock   float32 `json:"stores"`  // l库存数量
}
