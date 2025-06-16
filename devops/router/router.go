package router

import (
	"devops/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
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

	// API 路由组
	api := r.Group("/api")

	// 主机管理路由
	setupHostRoutes(api)

	// 仓库管理路由
	setupRepositoryRoutes(api)

	//镜像 中心
	RegisterDockerRegistryRoutes(api)

	//项目中心
	SetupProjectRoutes(api)

	return r
}

// setupHostRoutes 配置主机相关路由
func setupHostRoutes(api *gin.RouterGroup) {
	hostGroup := api.Group("/host")
	{
		hostGroup.GET("/list", controllers.GetHosts)
		hostGroup.POST("/add", controllers.CreateHost)
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
}

// setupRepositoryRoutes 配置仓库相关路由
func setupRepositoryRoutes(api *gin.RouterGroup) {
	repositoryController := controllers.NewRepositoryController()
	repos := api.Group("/repositories")
	{
		repos.POST("", repositoryController.CreateRepository)
		repos.GET("", repositoryController.GetRepositories)
		repos.PUT("/:id", repositoryController.UpdateRepository)
		repos.DELETE("/:id", repositoryController.DeleteRepository)
		repos.GET("/:id/branches", repositoryController.GetBranches)
		repos.GET("/:id/commits", repositoryController.GetCommits)
		repos.GET("/:id/files", repositoryController.GetFiles)
	}
}
