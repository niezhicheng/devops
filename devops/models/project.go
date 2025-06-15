package models

import (
	"gorm.io/gorm"
	"time"
)

// Project 项目模型
type Project struct {
	gorm.Model
	Name            string    `json:"name" gorm:"type:varchar(100);not null;comment:项目名称"`
	Description     string    `json:"description" gorm:"type:text;comment:项目描述"`
	RepositoryID    uint      `json:"repository_id" gorm:"not null;comment:关联的代码仓库ID"`
	Repository      Repository `json:"repository" gorm:"foreignKey:RepositoryID"`
	Branch          string    `json:"branch" gorm:"type:varchar(100);not null;comment:构建分支"`
	RegistryID      uint      `json:"registry_id" gorm:"not null;comment:关联的镜像仓库ID"`
	Registry        DockerRegistry `json:"registry" gorm:"foreignKey:RegistryID"`
	ImageName       string    `json:"image_name" gorm:"type:varchar(200);not null;comment:镜像名称"`
	ImageTag        string    `json:"image_tag" gorm:"type:varchar(100);not null;comment:镜像标签"`
	BuildScript     string    `json:"build_script" gorm:"type:text;comment:构建脚本"`
	Environment     string    `json:"environment" gorm:"type:varchar(50);comment:环境(dev/test/prod)"`
	Version         string    `json:"version" gorm:"type:varchar(50);comment:版本号"`
	BuildTimeout    int       `json:"build_timeout" gorm:"default:3600;comment:构建超时时间(秒)"`
	AutoBuild       bool      `json:"auto_build" gorm:"default:false;comment:是否自动构建"`
	BuildTriggers   string    `json:"build_triggers" gorm:"type:text;comment:构建触发器配置(JSON)"`
	LastBuildTime   time.Time `json:"last_build_time" gorm:"comment:最后构建时间"`
	LastBuildStatus string    `json:"last_build_status" gorm:"type:varchar(20);comment:最后构建状态"`
	CreatedBy       uint      `json:"created_by" gorm:"comment:创建人ID"`
	UpdatedBy       uint      `json:"updated_by" gorm:"comment:更新人ID"`
} 