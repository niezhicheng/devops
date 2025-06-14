package controllers

import (
	"devops/global"
	"devops/models"
	"devops/services"
	"devops/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var repositoryService = services.NewRepositoryService()

// List 获取仓库列表
func List(c *gin.Context) {
	var repositories []models.Repository
	if err := global.DB.Find(&repositories).Error; err != nil {
		utils.Logger.Error("获取仓库列表失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取仓库列表失败"})
		return
	}

	c.JSON(http.StatusOK, repositories)
}

// Create 创建仓库
func Create(c *gin.Context) {
	var repository models.Repository
	if err := c.ShouldBindJSON(&repository); err != nil {
		utils.Logger.Error("解析请求数据失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := global.DB.Create(&repository).Error; err != nil {
		utils.Logger.Error("创建仓库失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建仓库失败"})
		return
	}

	c.JSON(http.StatusCreated, repository)
}

// Get 获取仓库详情
func Get(c *gin.Context) {
	id := c.Param("id")
	var repository models.Repository

	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库详情失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	c.JSON(http.StatusOK, repository)
}

// Update 更新仓库
func Update(c *gin.Context) {
	id := c.Param("id")
	var repository models.Repository

	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	if err := c.ShouldBindJSON(&repository); err != nil {
		utils.Logger.Error("解析请求数据失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := global.DB.Save(&repository).Error; err != nil {
		utils.Logger.Error("更新仓库失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新仓库失败"})
		return
	}

	c.JSON(http.StatusOK, repository)
}

// Delete 删除仓库
func Delete(c *gin.Context) {
	id := c.Param("id")
	var repository models.Repository

	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	if err := global.DB.Delete(&repository).Error; err != nil {
		utils.Logger.Error("删除仓库失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除仓库失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "仓库已删除"})
}

// Sync 同步仓库
func Sync(c *gin.Context) {
	id := c.Param("id")
	var repository models.Repository

	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	if err := repositoryService.SyncRepository(&repository); err != nil {
		utils.Logger.Error("同步仓库失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步仓库失败"})
		return
	}

	c.JSON(http.StatusOK, repository)
}

// ListBranches 获取分支列表
func ListBranches(c *gin.Context) {
	id := c.Param("id")
	var repository models.Repository

	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	branches, err := repositoryService.GetBranches(&repository)
	if err != nil {
		utils.Logger.Error("获取分支列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分支列表失败"})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// ListFiles 获取文件列表
func ListFiles(c *gin.Context) {
	id := c.Param("id")
	path := c.DefaultQuery("path", "")
	branch := c.DefaultQuery("branch", "")
	fmt.Println(branch)

	var repository models.Repository
	if err := global.DB.First(&repository, id).Error; err != nil {
		utils.Logger.Error("获取仓库失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	files, err := repositoryService.GetFiles(&repository, path)
	if err != nil {
		utils.Logger.Error("获取文件列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败"})
		return
	}

	c.JSON(http.StatusOK, files)
}
