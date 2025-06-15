package controllers

import (
	"devops/global"
	"devops/models"
	"devops/services"
	"devops/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
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

// TestRepository 测试仓库连接
func TestRepository(ctx *gin.Context) {
	var req struct {
		URL   string `json:"url" binding:"required"`
		Token string `json:"token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	// 创建 GitHub 客户端
	client := github.NewClient(nil)
	baseURL, _ := url.Parse("https://api.github.com/")
	client.BaseURL = baseURL
	client.UserAgent = "devops"

	// 设置认证
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: req.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	// 测试连接 - 获取当前用户信息
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("GitHub API 连接失败: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接测试成功",
		"data": gin.H{
			"username": user.GetLogin(),
			"name":     user.GetName(),
		},
	})
}
