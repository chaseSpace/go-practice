package main

import (
	"github.com/gin-gonic/gin"
	"go_accost/api"
	"go_accost/config"
	"go_accost/crontask"
	"go_accost/db"
	"go_accost/handler"
	sugar "go_accost/log"
	"go_accost/util"
)

func main() {
	defer sugar.Stop()
	sugar.Info("config is ", util.Pretty(config.V))

	// 设置另一个服务的端点
	api.Init("http://1.1.1.1:8080/")

	db.Init()

	crontask.Start()
	defer crontask.Stop()

	eng := gin.Default()
	handler.Setup(eng)

	panic(eng.Run(":8080"))
}
