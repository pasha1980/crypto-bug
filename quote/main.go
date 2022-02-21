package quote

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote/config"
	"crypto-bug/quote/src/exchages"
	"crypto-bug/quote/src/service"
)

var Exchanges = []exchages.Exchange{
	exchages.Binance{},
	exchages.WhiteBit{},
	exchages.Coinlist{},
	exchages.ByBit{},
	exchages.Huobi{},
}

func Init() {
	var exception model.ExchangeException
	db := rootConfig.Database
	for _, exchange := range Exchanges {
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
	service.ClearQuotes()
}
