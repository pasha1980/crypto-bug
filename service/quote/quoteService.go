package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/telegram"
	"errors"
	"gorm.io/gorm"
	"os"
	"time"
)

const clearQuotesSql = `
DELETE 
FROM quotes 
WHERE date < ?
`

func ClearQuotes() {
	db := rootConfig.Database
	timeToLive, err := time.ParseDuration(os.Getenv("SAVE_QUOTE_TIME"))
	if err != nil {
		telegram.Log(err.Error(), "fatal")
	}

	if timeToLive == 0 {
		timeToLive = time.Hour
	}
	date := time.Now().Add(-timeToLive)
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

const getTrackCurrencyLastQuotesSQL = `
SELECT quotes.*
FROM
    quotes,
    (
        SELECT track_currency, MAX(date) AS max_date 
        FROM quotes 
        WHERE track_currency = ? 
        GROUP BY exchange, base_currency
    ) max_quotes
WHERE
    quotes.track_currency = max_quotes.track_currency
AND 
    quotes.date = max_quotes.max_date
GROUP BY id
`

func GetTrackCurrencyLastQuote(trackCurrency string) []model.Quote {
	db := rootConfig.Database
	var quotes []model.Quote
	err := db.
		Raw(getTrackCurrencyLastQuotesSQL, trackCurrency).
		Scan(&quotes).
		Error
	if err != nil {
		telegram.Log(err.Error(), "db")
	}
	return quotes
}

func GetMinMaxQuote(quotes []model.Quote) (error, *model.Quote, *model.Quote) {
	if len(quotes) == 0 {
		return errors.New("no quotes"), nil, nil
	}
	min := quotes[0]
	max := quotes[0]
	for _, quote := range quotes {
		if quote.Value > max.Value {
			max = quote
		}

		if quote.Value < min.Value {
			min = quote
		}
	}

	return nil, &min, &max
}

func IsSameCurrencyQuote(baseQuote model.Quote, trackQuote model.Quote) bool {
	return baseQuote.BaseCurrency == trackQuote.BaseCurrency && baseQuote.TrackCurrency == trackQuote.TrackCurrency
}
