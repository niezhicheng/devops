package models

import (
	"gorm.io/gorm"
	"time"
)

// Repository 仓库信息
type Repository struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	Name          string `gorm:"size:255;not null" json:"name"`
	Platform      string `gorm:"size:50;not null" json:"platform"`
	URL           string `gorm:"size:255;not null" json:"url"`
	Token         string `gorm:"size:255;not null" json:"token"`
	DefaultBranch string `gorm:"size:100" json:"defaultBranch"`
	Status        string `gorm:"size:50;not null;default:'active'" json:"status"`
	//LastSyncAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"lastSyncAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate 创建前钩子
//func (r *Repository) BeforeCreate(tx *gorm.DB) error {
//	if r.LastSyncAt.IsZero() {
//		r.LastSyncAt = time.Now()
//	}
//	return nil
//}

// TableName 设置表名
func (Repository) TableName() string {
	return "repositories"
}

// CreateRepository 创建仓库
func CreateRepository(db *gorm.DB, repo *Repository) error {
	return db.Create(repo).Error
}

// GetRepositoryList 获取仓库列表
func GetRepositoryList(db *gorm.DB, page, pageSize int, name, url string) ([]Repository, int64, error) {
	var repos []Repository
	var total int64

	query := db.Model(&Repository{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if url != "" {
		query = query.Where("url LIKE ?", "%"+url+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&repos).Error; err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}

// GetRepository 获取仓库详情
func GetRepository(db *gorm.DB, id uint) (*Repository, error) {
	var repo Repository
	if err := db.First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// UpdateRepository 更新仓库
func UpdateRepository(db *gorm.DB, id uint, repo *Repository) error {
	return db.Model(&Repository{}).Where("id = ?", id).Updates(repo).Error
}

// DeleteRepository 删除仓库
func DeleteRepository(db *gorm.DB, id uint) error {
	return db.Delete(&Repository{}, id).Error
}
