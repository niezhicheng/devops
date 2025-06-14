package controllers

import (
	"log"
	"net/http"
	"strconv"

	"devops/global"
	"devops/models"
	"github.com/gin-gonic/gin"
)

// CreateHost 添加主机
func CreateHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		log.Printf("绑定JSON失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateHost(global.DB, &host); err != nil {
		log.Printf("创建主机失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, host)
}

// GetHosts 获取主机列表
func GetHosts(c *gin.Context) {
	log.Printf("开始获取主机列表")

	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	name := c.Query("name")
	ip := c.Query("ip")

	log.Printf("查询参数: page=%d, pageSize=%d, name=%s, ip=%s", page, pageSize, name, ip)

	hosts, total, err := models.GetHostList(global.DB, page, pageSize, name, ip)
	if err != nil {
		log.Printf("获取主机列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("成功获取主机列表: 总数=%d", total)
	c.JSON(http.StatusOK, gin.H{
		"list":  hosts,
		"total": total,
	})
}

// UpdateHost 更新主机信息
func UpdateHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("无效的ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		log.Printf("绑定JSON失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateHost(global.DB, uint(id), &host); err != nil {
		log.Printf("更新主机失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, host)
}

// DeleteHost 删除主机
func DeleteHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("无效的ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteHost(global.DB, uint(id)); err != nil {
		log.Printf("删除主机失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Host deleted successfully"})
}
