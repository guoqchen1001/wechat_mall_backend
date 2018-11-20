package models

//User 用户
type User struct {
	UserID    string `gorm:"primary_key"` // 用户id
	UserNo    string // 用户编码
	UserName  string // 用户名
	WechatID  string `json:"openId"` // 微信id
	Phone     string // 手机号
	Country   string // 国家
	NickName  string `json:"nickName"`  //  昵称
	AvatarURL string `json:"avatarurl"` // 图像地址
	Gender    int    // 性别
	Province  string // 省份
	City      string // 城市
	Language  string // 显示语言

}
