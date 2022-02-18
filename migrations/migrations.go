package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"log"
)

func Migrate() {
	var err error
	db := config.Database
	err = db.AutoMigrate(&model.Quote{})
	err = db.AutoMigrate(&model.Statistic{})
	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		log.Fatal(err)
	}
}
