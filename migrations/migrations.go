package migrations

import (
	"crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
	"crypto-bug/quote/src/exchages"
)

func Migrate() {
	var err error
	db := config.Database
	err = db.AutoMigrate(&model.Quote{})
	err = db.AutoMigrate(&model.Statistic{})
	err = db.AutoMigrate(&model.ExchangeException{})
	ExchangeExceptionMigrate()

	// По мере добавления моделей добавлять новые миграции

	if err != nil {
		service.Log("Migration error: "+err.Error(), "migration")
	}
}

var ExchangeExceptions = []map[string]string{
	{
		"exchange": exchages.Coinlist{}.GetName(),
		"base":     "USDT",
		"track":    "ADA",
	},
	{
		"exchange": exchages.Coinlist{}.GetName(),
		"base":     "USDT",
		"track":    "BNB",
	},
}

func ExchangeExceptionMigrate() {
	var foundException model.ExchangeException
	db := config.Database
	for _, migration := range ExchangeExceptions {
		exception := model.ExchangeException{
			Exchange:      migration["exchange"],
			BaseCurrency:  migration["base"],
			TrackCurrency: migration["track"],
		}
		_ = db.Where(&exception).First(&foundException).Error
		if foundException.ID == 0 {
			db.Save(&exception)
		}

	}
}
