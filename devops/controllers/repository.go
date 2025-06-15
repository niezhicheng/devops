package controllers

import (
	"devops/global"
	"devops/models"
	"devops/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RepositoryController 仓库控制器
type RepositoryController struct {
	service *services.RepositoryService
}

// NewRepositoryController 创建仓库控制器
func NewRepositoryController() *RepositoryController {
	return &RepositoryController{
		service: services.NewRepositoryService(global.DB),
	}
}

// CreateRepository 创建仓库
func (c *RepositoryController) CreateRepository(ctx *gin.Context) {
	var repo models.Repository
	if err := ctx.ShouldBindJSON(&repo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateRepository(c.service.DB, &repo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repo)
}

// GetRepositories 获取仓库列表
func (c *RepositoryController) GetRepositories(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	name := ctx.Query("name")
	url := ctx.Query("url")

	repos, total, err := models.GetRepositoryList(c.service.DB, page, pageSize, name, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"list":  repos,
		"total": total,
	})
}

// UpdateRepository 更新仓库
func (c *RepositoryController) UpdateRepository(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var repo models.Repository
	if err := ctx.ShouldBindJSON(&repo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateRepository(c.service.DB, uint(id), &repo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repo)
}

// DeleteRepository 删除仓库
func (c *RepositoryController) DeleteRepository(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteRepository(c.service.DB, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Repository deleted successfully"})
}

// GetBranches 获取分支列表
func (c *RepositoryController) GetBranches(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	repo, err := models.GetRepository(c.service.DB, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	branches, err := c.service.GetBranches(ctx.Request.Context(), repo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, branches)
}

// GetCommits 获取提交历史
func (c *RepositoryController) GetCommits(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	branch := ctx.Query("branch")

	repo, err := models.GetRepository(c.service.DB, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commits, err := c.service.GetCommits(ctx.Request.Context(), repo, branch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, commits)
}

// GetFiles 获取文件列表
func (c *RepositoryController) GetFiles(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	path := ctx.Query("path")

	repo, err := models.GetRepository(c.service.DB, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files, err := c.service.GetFiles(repo, path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, files)
}
