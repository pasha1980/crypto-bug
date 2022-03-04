package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var Database *gorm.DB

func databaseSetup() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	Database = db
}
