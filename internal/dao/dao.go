package dao

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"kong-anime-go/internal/config"
)

// DB 数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	dbConfig := config.DBConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.Charset)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = Migrate(DB)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database connection and migration successful.")

	return DB
}
