package algorithm

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/quote"
	quote2 "crypto-bug/service/quote"
	"crypto-bug/service/telegram"
	"fmt"
	"math"
)

type ExchangeDifferenceAlgorithm struct {
}

const differenceThreshold = 5
const exchangeDifferenceBaseMessage = `
!РАЗНИЦА НА БИРЖАХ!
На бирже %s %s/%s - %g
На бирже %s %s/%s - %g
Разница - %.2f
`

func (algo ExchangeDifferenceAlgorithm) Analyze() {
	processedBaseExchange := make(map[string]bool)
	db := rootConfig.Database
	for _, baseExchange := range quote.Exchanges {
		baseQuotes := quote2.GetExchangeLastQuote(baseExchange.GetName())

		for _, trackExchange := range quote.Exchanges {

			if baseExchange.GetName() == trackExchange.GetName() {
				continue
			}

			if processedBaseExchange[trackExchange.GetName()] {
				continue
			}

			trackQuotes := quote2.GetExchangeLastQuote(trackExchange.GetName())

			for _, baseQuote := range baseQuotes {
				for _, trackQuote := range trackQuotes {
					if !quote2.IsSameCurrencyQuote(baseQuote, trackQuote) {
						continue
					}

					diff := algo.CalculateDiff(baseQuote.Value, trackQuote.Value)
					if diff >= differenceThreshold {

						message := fmt.Sprintf(
							exchangeDifferenceBaseMessage,
							baseExchange.GetName(),
							baseQuote.TrackCurrency,
							baseQuote.BaseCurrency,
							baseQuote.Value,
							trackExchange.GetName(),
							trackQuote.TrackCurrency,
							trackQuote.BaseCurrency,
							trackQuote.Value,
							diff,
						)
						tgMessage := telegram.TelegramMessage{
							Text: message,
						}
						tgMessage.SendToMe()

						statistic := model.Statistic{
							AlgorithmName: algo.GetName(),
							BaseCurrency:  trackQuote.BaseCurrency,
							TrackCurrency: trackQuote.TrackCurrency,
							Exchange:      baseExchange.GetName() + "/" + trackExchange.GetName(),
							Result:        fmt.Sprintf("%f", diff),
							Action:        "Sent telegram message",
						}
						db.Save(&statistic)
					}
				}
			}
		}
		processedBaseExchange[baseExchange.GetName()] = true
	}
}

func (algo ExchangeDifferenceAlgorithm) CalculateDiff(base float64, track float64) float64 {
	return math.Abs(((track - base) / base) * 100)
}

func (algo ExchangeDifferenceAlgorithm) GetName() string {
	return "Exchange difference"
}
