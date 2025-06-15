package main

import (
	"devops/config"
	"devops/router"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 配置路由
	r := router.SetupRouter()

	// 启动服务器
	r.Run(":8080")
}
