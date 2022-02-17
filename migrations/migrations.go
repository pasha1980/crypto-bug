package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"log"
)

func Migrate() {
	db := config.Database
	err := db.AutoMigrate(&model.Quote{})
	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		log.Fatal(err)
	}
}
