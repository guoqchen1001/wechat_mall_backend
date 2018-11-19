package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const logName = "wechat.log"

func init() {

	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("打开或创建日志文件失败，将只用标准输入输出流来操作日志")
	} else {
		log.Out = file
	}

	log.WithFields(logrus.Fields{
		"start": "wechat_mall",
	}).Info("程序启动")

}
