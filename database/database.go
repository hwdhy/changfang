package database

import (
	"changfang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	//连接数据库
	dsn := "host=localhost user=postgres password=dhyy dbname=changfang port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
		return nil
	}
	//自动创建表
	db.AutoMigrate(
		&ServiceModels.QueueUrl{},
		&ServiceModels.UrlContent{},
	)

	return db

}
