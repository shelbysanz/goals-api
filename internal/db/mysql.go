package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	// data source name
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		"g_user",
		"TestPass123!",
		"127.0.0.1",
		"3306",
		"goals",
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
