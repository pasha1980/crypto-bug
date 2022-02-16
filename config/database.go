package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Database *gorm.DB

func databaseSetup()  {
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Database = db
}

