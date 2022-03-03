package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote/config"
	"crypto-bug/service/telegram"
	"errors"
	"gorm.io/gorm"
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
	err := db.Where(&exception).First(&foundException).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.Save(&exception)
	}
}

const getExchangeLastQuotesSQL = `
SELECT quotes.* 
FROM 
    quotes, 
    (SELECT exchange, MAX(date) AS max_date FROM quotes WHERE exchange = ? GROUP BY track_currency, base_currency) max_quotes
WHERE 
    quotes.exchange = max_quotes.exchange 
AND 
    quotes.date = max_quotes.max_date
`

func GetExchangeLastQuote(exchange string) []model.Quote {
	db := rootConfig.Database
	var quotes []model.Quote
	err := db.
		Raw(getExchangeLastQuotesSQL, exchange).
		Scan(&quotes).
		Error
	if err != nil {
		telegram.Log(err.Error(), "db")
	}
	return quotes
}

func IsSameCurrencyQuote(baseQuote model.Quote, trackQuote model.Quote) bool {
	return baseQuote.BaseCurrency == trackQuote.BaseCurrency && baseQuote.TrackCurrency == trackQuote.TrackCurrency
}
