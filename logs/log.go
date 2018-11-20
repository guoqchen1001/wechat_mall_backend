package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

//Log 日志对象
var Log = logrus.New()

const logName = "wechat.log"

func init() {

	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Info("打开或创建日志文件失败，将只用标准输入输出流来操作日志")
	} else {
		Log.Out = file
	}

	Log.WithFields(logrus.Fields{
		"start": "wechat_mall",
	}).Info("程序启动")

}
