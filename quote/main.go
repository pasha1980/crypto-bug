package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quoteConf"
	"crypto-bug/service/quote"
	"errors"
	"gorm.io/gorm"
)

func Init() {
	var err error
	db := rootConfig.Database
	for _, exchange := range quoteConf.Exchanges {
		for _, trackCurrency := range quoteConf.CurrenciesToTrack {
			for _, baseCurrency := range quoteConf.BaseCurrencies {
				var exception model.ExchangeException
				err = db.Where(&model.ExchangeException{
					Exchange:      exchange.GetName(),
					TrackCurrency: trackCurrency,
					BaseCurrency:  baseCurrency,
				}).First(&exception).Error
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				exchange.Save(trackCurrency, baseCurrency)
			}
		}
	}
	quote.ClearQuotes()
}
