package config

import (
	"fmt"
	"log"
	"time"

	"devops/global"
	"devops/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	dsn := "root:d200145001@tcp(127.0.0.1:3306)/devops?charset=utf8mb4&parseTime=True&loc=Local"

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := global.DB.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 自动迁移数据库表
	log.Println("开始数据库迁移...")
	err = global.DB.AutoMigrate(&models.Host{}, models.Repository{}, models.DockerRegistry{}, models.Project{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	fmt.Println("数据库连接成功")
}
