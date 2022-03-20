package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/chain"
	"crypto-bug/service/telegram"
)

func Migrate() {
	DatabaseMigration()
	CacheMigration()
}

func DatabaseMigration() {
	var err error
	db := config.Database
	err = db.AutoMigrate(&model.Quote{})
	err = db.AutoMigrate(&model.Statistic{})
	err = db.AutoMigrate(&model.ExchangeException{})

	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		telegram.Log("Database migration error: "+err.Error(), "migration")
	}
}

func CacheMigration() {
	var err error
	err = chain.Bep20{}.GetTokens()

	// По мере необходимости добавлять новые миграции

	if err != nil {
		telegram.Log("Cache migration error: "+err.Error(), "migration")
	}
}
