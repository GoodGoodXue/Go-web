package main

import (
	"Practice/Go-web/Go_web_v2/common"
	"Practice/Go-web/Go_web_v2/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// 调用连接数据库，启动
	db := common.InitDB()

	// 返回前关闭
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	engin := gin.Default()

	routes.CollectRoutes(engin)

	panic(engin.Run(":9090"))
}
