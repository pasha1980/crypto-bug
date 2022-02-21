package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote/config"
	"crypto-bug/quote/src/exchages"
	"time"
)

func Init() {
	var exception model.ExchangeException
	db := rootConfig.Database
	for _, exchange := range config.Exchanges {
		exchangeObj := exchange
		for _, trackCurrency := range config.CurrenciesToTrack {
			for _, baseCurrency := range config.BaseCurrencies {
				_ = db.Where(model.ExchangeException{
					Exchange:      exchange.GetName(),
					TrackCurrency: trackCurrency,
					BaseCurrency:  baseCurrency,
				}).First(&exception).Error
				if exception.ID != 0 {
					continue
				}
				exchangeObj.Save(trackCurrency, baseCurrency)
			}
		}
	}
	ClearQuotes()
}

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

func ProcessException(exchange exchages.Exchange, track string, base string) {
	var foundException model.ExchangeException
	exception := model.ExchangeException{
		Exchange:      exchange.GetName(),
		TrackCurrency: track,
		BaseCurrency:  base,
	}

	db := rootConfig.Database
	_ = db.Where(&exception).First(&foundException).Error
	if foundException.ID == 0 {
		db.Save(&exception)
	}
}
