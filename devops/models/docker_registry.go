package models

import (
	"gorm.io/gorm"
	"time"
)

// DockerRegistry Docker 镜像仓库
type DockerRegistry struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:100;not null;comment:仓库名称"`
	Type      string    `json:"type" gorm:"size:50;not null;comment:仓库类型(public/private)"`
	URL       string    `json:"url" gorm:"size:255;not null;comment:仓库地址"`
	Username  string    `json:"username" gorm:"size:100;comment:用户名"`
	Password  string    `json:"password" gorm:"size:255;comment:密码"`
	Status    string    `json:"status" gorm:"size:20;default:'active';comment:状态"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateDockerRegistry 创建 Docker 镜像仓库
func CreateDockerRegistry(db *gorm.DB, registry *DockerRegistry) error {
	return db.Create(registry).Error
}

// GetDockerRegistry 获取 Docker 镜像仓库
func GetDockerRegistry(db *gorm.DB, id uint) (*DockerRegistry, error) {
	var registry DockerRegistry
	err := db.First(&registry, id).Error
	return &registry, err
}

// GetDockerRegistryList 获取 Docker 镜像仓库列表
func GetDockerRegistryList(db *gorm.DB, page, pageSize int, name string) ([]DockerRegistry, int64, error) {
	var registries []DockerRegistry
	var total int64

	query := db.Model(&DockerRegistry{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&registries).Error
	return registries, total, err
}

// UpdateDockerRegistry 更新 Docker 镜像仓库
func UpdateDockerRegistry(db *gorm.DB, id uint, registry *DockerRegistry) error {
	return db.Model(&DockerRegistry{}).Where("id = ?", id).Updates(registry).Error
}

// DeleteDockerRegistry 删除 Docker 镜像仓库
func DeleteDockerRegistry(db *gorm.DB, id uint) error {
	return db.Delete(&DockerRegistry{}, id).Error
} 