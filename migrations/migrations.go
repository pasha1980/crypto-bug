package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/parserService"
)

func Migrate() {
	var err error
	db := config.Database
	err = db.AutoMigrate(&model.Quote{})
	err = db.AutoMigrate(&model.Statistic{})
	err = db.AutoMigrate(&model.ExchangeException{})

	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		parserService.Log("Migration error: "+err.Error(), "migration")
	}
}
