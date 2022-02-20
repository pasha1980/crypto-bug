package algorithm

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
	quoteConfig "crypto-bug/quote/config"
	"fmt"
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
					service.Log("Error getting quotes from db", "algo")
					continue
				}

				err = query.Not("id = ?", lastQuote.ID).Order("id desc").Find(&quotes).Error
				if err != nil {
					service.Log("Error getting quotes from db", "algo")
					continue
				}

				for _, baseQuote := range quotes {
					diff = algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
					if diff > 0 && diff >= growthThreshold {
						algo.AbnormalGrowth(lastQuote, baseQuote)
						break
					} else if diff < 0 && -diff >= dropThreshold {
						algo.AbnormalDrop(lastQuote, baseQuote)
						break
					}
				}

			}
		}
	}
}

func (algo AbnormallyPriceAlgorithm) GetName() string {
	return "Abnormally price"
}

const baseAbnormalGrowthMessage = `
На бирже %s аномальный рост %s/%s
%s цена составляла %g %s
Последняя котировка - %g %s
Разница - %.2f
`

func (algo AbnormallyPriceAlgorithm) AbnormalGrowth(lastQuote model.Quote, baseQuote model.Quote) {
	diff := algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
	message := service.TelegramMessage{
		Text: fmt.Sprintf(
			baseAbnormalGrowthMessage,
			lastQuote.Exchange,
			lastQuote.TrackingCurrency,
			lastQuote.BaseCurrency,
			baseQuote.Date.Format("02.01.2006 15:04:05"),
			baseQuote.Value,
			baseQuote.BaseCurrency,
			lastQuote.Value,
			lastQuote.BaseCurrency,
			diff,
		),
	}
	message.SendToMe()

	db := rootConfig.Database
	statistic := model.Statistic{
		AlgorithmName: algo.GetName(),
		BaseCurrency:  lastQuote.BaseCurrency,
		TrackCurrency: lastQuote.TrackingCurrency,
		Exchange:      lastQuote.Exchange,
		Result:        fmt.Sprintf("%f", diff),
		Action:        "Sent telegram message",
	}
	db.Save(&statistic)
}

const baseAbnormalDropMessage = `
На бирже %s аномальное падение %s/%s
%s цена составляла %g %s
Последняя котировка - %g %s
Разница - %.2f
`

func (algo AbnormallyPriceAlgorithm) AbnormalDrop(lastQuote model.Quote, baseQuote model.Quote) {
	diff := algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
	message := service.TelegramMessage{
		Text: fmt.Sprintf(
			baseAbnormalDropMessage,
			lastQuote.Exchange,
			lastQuote.TrackingCurrency,
			lastQuote.BaseCurrency,
			baseQuote.Date.Format("02.01.2006 15:04:05"),
			baseQuote.Value,
			baseQuote.BaseCurrency,
			lastQuote.Value,
			lastQuote.BaseCurrency,
			-diff,
		),
	}
	message.SendToMe()

	db := rootConfig.Database
	statistic := model.Statistic{
		AlgorithmName: algo.GetName(),
		BaseCurrency:  lastQuote.BaseCurrency,
		TrackCurrency: lastQuote.TrackingCurrency,
		Exchange:      lastQuote.Exchange,
		Result:        fmt.Sprintf("%f", diff),
		Action:        "Sent telegram message",
	}
	db.Save(&statistic)
}

func (algo AbnormallyPriceAlgorithm) CalculateDiff(lastQuote float64, baseQuote float64) float64 {
	return ((lastQuote - baseQuote) / baseQuote) * 100
}
