// wechat_mall_backend project main.go
package main

import (
	"net/http"
	"os"
	"path/filepath"
	"wechat_mall_backend/logs"

	"github.com/sirupsen/logrus"
)

func main() {

	// 创建路由
	err := CeateRouter()
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"webserver": "error",
		}).Error(err)
	}

	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"webserver": "error",
		}).Error(err)
	}

	certPem := filepath.Join(dir, "pem/cert.pem")
	keyPem := filepath.Join(dir, "pem/key.pem")

	err = server.ListenAndServeTLS(certPem, keyPem)

	logs.Log.WithFields(logrus.Fields{
		"webserver": "error",
	}).Error(err)

}
