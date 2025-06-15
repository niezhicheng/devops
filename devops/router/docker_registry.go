package router

import (
	"devops/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterDockerRegistryRoutes 注册 Docker 镜像仓库路由
func RegisterDockerRegistryRoutes(r *gin.RouterGroup) {
	registryController := controllers.NewDockerRegistryController()
	registry := r.Group("/docker-registries")
	{
		registry.POST("", registryController.CreateDockerRegistry)
		registry.GET("", registryController.GetDockerRegistries)
		registry.PUT("/:id", registryController.UpdateDockerRegistry)
		registry.DELETE("/:id", registryController.DeleteDockerRegistry)
		registry.POST("/test-connection", registryController.TestDockerRegistryConnection)
	}
} 