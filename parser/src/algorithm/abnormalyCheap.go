package algorithm

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	quoteConfig "crypto-bug/quote/config"
	"log"
)

type AbnormallyPriceAlgorithm struct {
}

const growthThreshold = 30.0
const dropThreshold = 40.0

func (algo AbnormallyPriceAlgorithm) Analyze() {
	var quotes []model.Quote
	var lastQuote model.Quote
	var diff float64
	var err error

	db := rootConfig.Database

	for _, exchange := range quoteConfig.Exchanges {
		for _, trackCurrency := range quoteConfig.CurrenciesToTrack {
			for _, baseCurrency := range quoteConfig.BaseCurrencies {
				query := db.Where(&model.Quote{
					Exchange:         exchange.GetName(),
					BaseCurrency:     baseCurrency,
					TrackingCurrency: trackCurrency,
				})

				err = db.Where(&model.Quote{
					Exchange:         exchange.GetName(),
					BaseCurrency:     baseCurrency,
					TrackingCurrency: trackCurrency,
				}).Last(&lastQuote).Error
				if err != nil {
					log.Println("Error getting quotes from db")
					continue
				}

				err = query.Find(&quotes).Error
				if err != nil {
					log.Println("Error getting quotes from db")
					continue
				}

				for _, baseQuote := range quotes {
					diff = algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
					if diff > 0 && diff >= growthThreshold {
						algo.AbnormalGrowth(lastQuote, baseQuote)
					} else if diff < 0 && -diff >= dropThreshold {
						algo.AbnormalDrop(lastQuote, baseQuote)
					}
				}

			}
		}
	}
}

func (algo AbnormallyPriceAlgorithm) AbnormalGrowth(lastQuote model.Quote, baseQuote model.Quote) {

}

func (algo AbnormallyPriceAlgorithm) AbnormalDrop(lastQuote model.Quote, baseQuote model.Quote) {

}

func (algo AbnormallyPriceAlgorithm) CalculateDiff(lastQuote float64, baseQuote float64) float64 {
	return (lastQuote - baseQuote) / baseQuote
}
