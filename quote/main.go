package quote

import (
	config2 "crypto-bug/config"
	"crypto-bug/quote/config"
	"time"
)

func Init() {
	for _, exchange := range config.Exchanges {
		exchangeObj := exchange
		for _, trackCurrency := range config.CurrenciesToTrack {
			for _, baseCurrency := range config.BaseCurrencies {
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
	db := config2.Database
	date := time.Now().Add(-(time.Hour * config.HoursToSaveQuotes))
	db.Exec(clearQuotesSql, date)
}
