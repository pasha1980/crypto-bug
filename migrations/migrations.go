package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/telegram"
)

func Migrate() {
	var err error
	db := config.Database
	err = db.AutoMigrate(&model.Quote{})
	err = db.AutoMigrate(&model.Statistic{})
	err = db.AutoMigrate(&model.ExchangeException{})

	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		telegram.Log("Migration error: "+err.Error(), "migration")
	}
}
