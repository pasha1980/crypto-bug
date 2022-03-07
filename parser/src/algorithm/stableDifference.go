package algorithm

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	quoteConfig "crypto-bug/quote/config"
	quoteService "crypto-bug/service/quote"
	"crypto-bug/service/telegram"
	"fmt"
)

type StableDifferenceAlgorithm struct {
}

const stableDifferenceThreshold = 5
const stableDifferenceMessage = `
РАЗНИЦА СТЕЙБЛОВ

На бирже %s %s/%s - %g
На бирже %s %s/%s - %g
Разница - %.2f
`

func (algo StableDifferenceAlgorithm) Analyze() {
	db := rootConfig.Database
	for _, trackCurrency := range quoteConfig.CurrenciesToTrack {
		err, minQuote, maxQuote := quoteService.GetMinMaxQuote(
			quoteService.GetTrackCurrencyLastQuote(trackCurrency),
		)
		if err != nil {
			continue
		}

		diff := algo.CalculateDiff(minQuote.Value, maxQuote.Value)
		if diff >= stableDifferenceThreshold {
			message := fmt.Sprintf(
				stableDifferenceMessage,
				minQuote.Exchange,
				minQuote.TrackCurrency,
				minQuote.BaseCurrency,
				minQuote.Value,
				maxQuote.Exchange,
				maxQuote.TrackCurrency,
				maxQuote.BaseCurrency,
				maxQuote.Value,
				diff,
			)
			tgMessage := telegram.TelegramMessage{
				Text: message,
			}
			tgMessage.SendToMe()

			statistic := model.Statistic{
				AlgorithmName: algo.GetName(),
				BaseCurrency:  minQuote.BaseCurrency + "/" + maxQuote.BaseCurrency,
				TrackCurrency: trackCurrency,
				Exchange:      minQuote.Exchange + "/" + minQuote.Exchange,
				Result:        fmt.Sprintf("%f", diff),
				Action:        "Sent telegram message",
			}
			db.Save(&statistic)
		}

	}
}

func (algo StableDifferenceAlgorithm) CalculateDiff(minQuote float64, maxQuote float64) float64 {
	return ((maxQuote - minQuote) / maxQuote) * 100
}

func (algo StableDifferenceAlgorithm) GetName() string {
	return "Stable Difference"
}
