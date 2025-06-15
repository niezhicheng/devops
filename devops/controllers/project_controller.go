package controllers

import (
	"devops/global"
	"devops/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ProjectController 项目控制器
type ProjectController struct {
	DB *gorm.DB
}

// NewProjectController 创建项目控制器
func NewProjectController() *ProjectController {
	return &ProjectController{
		DB: global.DB,
	}
}

// CreateProject 创建项目
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置创建人
	project.CreatedBy = ctx.GetUint("user_id")
	project.UpdatedBy = ctx.GetUint("user_id")

	if err := c.DB.Create(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "项目创建成功",
		"data":    project,
	})
}

// GetProjects 获取项目列表
func (c *ProjectController) GetProjects(ctx *gin.Context) {
	var projects []models.Project
	query := c.DB.Model(&models.Project{})

	// 添加过滤条件
	if name := ctx.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if environment := ctx.Query("environment"); environment != "" {
		query = query.Where("environment = ?", environment)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("Repository").Preload("Registry").
		Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": projects,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProject 获取项目详情
func (c *ProjectController) GetProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project

	if err := c.DB.Preload("Repository").Preload("Registry").
		First(&project, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

// UpdateProject 更新项目
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project

	if err := c.DB.First(&project, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置更新人
	project.UpdatedBy = ctx.GetUint("user_id")

	if err := c.DB.Save(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "项目更新成功",
		"data":    project,
	})
}

// DeleteProject 删除项目
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Project{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "项目删除成功"})
} 