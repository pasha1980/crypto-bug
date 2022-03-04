package algorithm

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote"
	quoteConfig "crypto-bug/quote/config"
	"crypto-bug/service/telegram"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AbnormallyPriceAlgorithm struct {
}

const growthThreshold = 20.0
const dropThreshold = 20.0

func (algo AbnormallyPriceAlgorithm) Analyze() {
	var quotes []model.Quote
	var lastQuote model.Quote
	var diff float64
	var err error

	db := rootConfig.Database

	for _, exchange := range quote.Exchanges {
		for _, trackCurrency := range quoteConfig.CurrenciesToTrack {
			for _, baseCurrency := range quoteConfig.BaseCurrencies {
				query := db.Where(&model.Quote{
					Exchange:      exchange.GetName(),
					BaseCurrency:  baseCurrency,
					TrackCurrency: trackCurrency,
					IsAbnormally:  false,
				})

				err = db.Where(&model.Quote{
					Exchange:      exchange.GetName(),
					BaseCurrency:  baseCurrency,
					TrackCurrency: trackCurrency,
				}).Last(&lastQuote).Error
				if err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						telegram.Log("Error getting quotes from db", "algo")
					}
					continue
				}

				err = query.Not("id = ?", lastQuote.ID).Order("id desc").Find(&quotes).Error
				if err != nil {
					telegram.Log("Error getting quotes from db", "algo")
					continue
				}

				for _, baseQuote := range quotes {
					diff = algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
					if diff > 0 && diff >= growthThreshold {
						algo.ActionAbnormalGrowth(lastQuote, baseQuote)
						break
					} else if diff < 0 && -diff >= dropThreshold {
						algo.ActionAbnormalDrop(lastQuote, baseQuote)
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
**АНОМАЛЬНЫЙ РОСТ**
На бирже %s аномальный рост %s/%s
%s цена составляла %g %s
Последняя котировка - %g %s
Разница - %.2f
`

func (algo AbnormallyPriceAlgorithm) ActionAbnormalGrowth(lastQuote model.Quote, baseQuote model.Quote) {
	diff := algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
	message := telegram.TelegramMessage{
		Text: fmt.Sprintf(
			baseAbnormalGrowthMessage,
			lastQuote.Exchange,
			lastQuote.TrackCurrency,
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
		TrackCurrency: lastQuote.TrackCurrency,
		Exchange:      lastQuote.Exchange,
		Result:        fmt.Sprintf("%f", diff),
		Action:        "Sent telegram message",
	}
	db.Save(&statistic)

	lastQuote.IsAbnormally = true
	db.Save(&lastQuote)
}

const baseAbnormalDropMessage = `
**АНОМАЛЬНОЕ ПАДЕНИЕ**
На бирже %s аномальное падение %s/%s
%s цена составляла %g %s
Последняя котировка - %g %s
Разница - %.2f
`

func (algo AbnormallyPriceAlgorithm) ActionAbnormalDrop(lastQuote model.Quote, baseQuote model.Quote) {
	diff := algo.CalculateDiff(lastQuote.Value, baseQuote.Value)
	message := telegram.TelegramMessage{
		Text: fmt.Sprintf(
			baseAbnormalDropMessage,
			lastQuote.Exchange,
			lastQuote.TrackCurrency,
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
		TrackCurrency: lastQuote.TrackCurrency,
		Exchange:      lastQuote.Exchange,
		Result:        fmt.Sprintf("%f", diff),
		Action:        "Sent telegram message",
	}
	db.Save(&statistic)

	lastQuote.IsAbnormally = true
	db.Save(&lastQuote)
}

func (algo AbnormallyPriceAlgorithm) CalculateDiff(lastQuote float64, baseQuote float64) float64 {
	return ((lastQuote - baseQuote) / baseQuote) * 100
}
