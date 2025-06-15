package router

import (
	"devops/controllers"
	"github.com/gin-gonic/gin"
)

// SetupProjectRoutes 设置项目路由
func SetupProjectRoutes(router *gin.RouterGroup) {
	projectController := controllers.NewProjectController()

	// 项目路由组
	projects := router.Group("/projects")
	{
		projects.POST("", projectController.CreateProject)
		projects.GET("", projectController.GetProjects)
		projects.GET("/:id", projectController.GetProject)
		projects.PUT("/:id", projectController.UpdateProject)
		projects.DELETE("/:id", projectController.DeleteProject)
	}
}
