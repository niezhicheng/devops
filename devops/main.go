package main

import (
	"devops/config"
	"devops/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 创建 Gin 路由
	r := gin.Default()

	// 设置 CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 主机管理路由
	hostGroup := r.Group("/api/host")
	{
		hostGroup.GET("/list", controllers.GetHosts)
		hostGroup.POST("", controllers.CreateHost)
		hostGroup.PUT("/:id", controllers.UpdateHost)
		hostGroup.DELETE("/:id", controllers.DeleteHost)

		// SFTP相关路由
		hostGroup.GET("/:id/sftp", controllers.GetSftpFiles)
		hostGroup.POST("/:id/sftp/upload", controllers.UploadSftpFile)
		hostGroup.GET("/:id/sftp/download", controllers.DownloadSftpFile)
		hostGroup.GET("/:id/sftp/download-dir", controllers.DownloadSftpDir)
		hostGroup.DELETE("/:id/sftp", controllers.DeleteSftpFile)
		hostGroup.PUT("/:id/sftp/rename", controllers.RenameSftpFile)
		hostGroup.POST("/:id/sftp/compress", controllers.CompressSftpDir)
		hostGroup.GET("/:id/webshell", controllers.WebShell)
		hostGroup.POST("/:id/upload", controllers.UploadFile)
		hostGroup.GET("/:id/download", controllers.DownloadFile)
	}

	// 仓库管理路由
	repositoryGroup := r.Group("/api/repositories")
	//repositoryGroup.Use(middleware.AuthMiddleware())
	{
		repositoryGroup.GET("", controllers.List)
		repositoryGroup.POST("", controllers.Create)
		repositoryGroup.GET("/:id", controllers.Get)
		repositoryGroup.PUT("/:id", controllers.Update)
		repositoryGroup.DELETE("/:id", controllers.Delete)
	}

	// 启动服务器
	r.Run(":8080")
}
