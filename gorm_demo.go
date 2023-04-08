package utils

import "gorm.io/gorm"

var db *gorm.DB

// GetDB 获取连接
func GetDB() *gorm.DB {
	return db
}

// initializes the database instance
func Init() {
	dbConf := &DatabaseConfig{
		DbMode:      true,
		DbName:      "db_name",
		Type:        "mysql",
		Dsn:         "user:password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 3600,
	}
	db = NewInstance(dbConf, &gorm.Config{})
}

// CloseDB 关闭连接
func CloseDB() {
	defer ReleaseInstance()
}
