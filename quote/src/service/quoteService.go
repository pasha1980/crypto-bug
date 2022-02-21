package service

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote/config"
	"time"
)

const clearQuotesSql = `
DELETE 
FROM quotes 
WHERE date < ?
`

func ClearQuotes() {
	db := rootConfig.Database
	date := time.Now().Add(-(time.Hour * config.HoursToSaveQuotes))
	db.Exec(clearQuotesSql, date)
}

func ProcessException(exchange string, track string, base string) {
	var foundException model.ExchangeException
	exception := model.ExchangeException{
		Exchange:      exchange,
		TrackCurrency: track,
		BaseCurrency:  base,
	}

	db := rootConfig.Database
	_ = db.Where(&exception).First(&foundException).Error
	if foundException.ID == 0 {
		db.Save(&exception)
	}
}
