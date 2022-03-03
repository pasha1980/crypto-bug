package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote/config"
	"crypto-bug/quote/src/exchages"
	"crypto-bug/service/quote"
	"errors"
	"gorm.io/gorm"
)

var Exchanges = []exchages.Exchange{
	exchages.Binance{},
	exchages.WhiteBit{},
	exchages.Coinlist{},
	exchages.ByBit{},
	exchages.Huobi{},
}

func Init() {
	var err error
	db := rootConfig.Database
	for _, exchange := range Exchanges {
		for _, trackCurrency := range config.CurrenciesToTrack {
			for _, baseCurrency := range config.BaseCurrencies {
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
