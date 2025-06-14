package models

import (
	"gorm.io/gorm"
	"time"
)

// Host 主机模型
type Host struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	IP          string    `gorm:"size:15;not null" json:"ip"`
	Port        int       `gorm:"not null" json:"port"`
	Username    string    `gorm:"size:50;not null" json:"username"`
	Password    string    `gorm:"size:100;not null" json:"password"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"createdTime"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" json:"updatedTime"`
}

// TableName 指定表名
func (Host) TableName() string {
	return "hosts"
}

// CreateHost 创建主机
func CreateHost(db *gorm.DB, host *Host) error {
	return db.Create(host).Error
}

// GetHostList 获取主机列表
func GetHostList(db *gorm.DB, page, pageSize int, name, ip string) ([]Host, int64, error) {
	var hosts []Host
	var total int64

	query := db.Model(&Host{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&hosts).Error
	return hosts, total, err
}

// UpdateHost 更新主机信息
func UpdateHost(db *gorm.DB, id uint, host *Host) error {
	return db.Model(&Host{}).Where("id = ?", id).Updates(host).Error
}

// DeleteHost 删除主机
func DeleteHost(db *gorm.DB, id uint) error {
	return db.Delete(&Host{}, id).Error
}

// GetHostByID 根据ID获取主机
func GetHostByID(db *gorm.DB, id uint) (*Host, error) {
	var host Host
	err := db.First(&host, id).Error
	return &host, err
}
