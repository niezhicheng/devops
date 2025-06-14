package models

import (
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
