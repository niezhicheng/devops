package controllers

import (
	"context"
	"devops/global"
	"devops/models"
	"net/http"
	"strconv"
	"strings"

	dockerregistry "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DockerRegistryController Docker 镜像仓库控制器
type DockerRegistryController struct {
	DB *gorm.DB
}

// NewDockerRegistryController 创建 Docker 镜像仓库控制器
func NewDockerRegistryController() *DockerRegistryController {
	return &DockerRegistryController{
		DB: global.DB,
	}
}

// CreateDockerRegistry 创建 Docker 镜像仓库
func (c *DockerRegistryController) CreateDockerRegistry(ctx *gin.Context) {
	var registry models.DockerRegistry
	if err := ctx.ShouldBindJSON(&registry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateDockerRegistry(c.DB, &registry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registry)
}

// GetDockerRegistries 获取 Docker 镜像仓库列表
func (c *DockerRegistryController) GetDockerRegistries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	name := ctx.Query("name")

	registries, total, err := models.GetDockerRegistryList(c.DB, page, pageSize, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"list":  registries,
		"total": total,
	})
}

// UpdateDockerRegistry 更新 Docker 镜像仓库
func (c *DockerRegistryController) UpdateDockerRegistry(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var registry models.DockerRegistry
	if err := ctx.ShouldBindJSON(&registry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateDockerRegistry(c.DB, uint(id), &registry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registry)
}

// DeleteDockerRegistry 删除 Docker 镜像仓库
func (c *DockerRegistryController) DeleteDockerRegistry(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteDockerRegistry(c.DB, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Docker registry deleted successfully"})
}

// TestDockerRegistryConnection 测试 Docker 镜像仓库连接
func (c *DockerRegistryController) TestDockerRegistryConnection(ctx *gin.Context) {
	var registry models.DockerRegistry
	if err := ctx.ShouldBindJSON(&registry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保 URL 不包含协议前缀
	url := registry.URL
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// 创建 Docker 客户端，使用自动检测的 API 版本
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(), // 自动协商 API 版本
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "无法创建 Docker 客户端: " + err.Error(),
			"status": "error",
		})
		return
	}
	defer cli.Close()

	// 设置认证信息
	authConfig := dockerregistry.AuthConfig{
		Username:      registry.Username,
		Password:      registry.Password,
		ServerAddress: url,
	}

	// 尝试登录
	_, err = cli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":  "认证失败: " + err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "连接成功",
		"status":  "success",
	})
}
