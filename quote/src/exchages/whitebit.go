package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/quote"
	"crypto-bug/service/telegram"
	"encoding/json"
	"strconv"
	"time"
)

type WhiteBit struct {
}

type WhiteBitResponse struct {
	Success bool                   `json:"success"`
	Message WhiteBitMessage        `json:"message"`
	Result  []WhiteBitMarketResult `json:"result"`
}

type WhiteBitMessage struct {
	Market []string `json:"market"`
}

type WhiteBitMarketResult struct {
	Price string `json:"price"`
}

const whiteBitReturnMessageNeedException = "Market is not available."

func (whiteBit WhiteBit) Save(track string, base string) {
	var response WhiteBitResponse
	market := track + "_" + base
	client := rootConfig.Client

	responseRaw, err := client.Get("https://whitebit.com/api/v1/public/history?lastId=1&limit=1&market=" + market)
	if err != nil {
		telegram.Log("WhiteBit connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if !response.Success {
		if len(response.Message.Market) == 0 {
			return
		} else if response.Message.Market[0] == whiteBitReturnMessageNeedException {
			quote.ProcessException(whiteBit.GetName(), track, base)
		} else {
			telegram.Log("WhiteBit request error. "+response.Message.Market[0], "exchange")
		}
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Result[0].Price, 64)

	newQuote := model.Quote{
		Exchange:      whiteBit.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         priceFloat,
	}

	db := rootConfig.Database
	db.Save(&newQuote)
}

func (whiteBit WhiteBit) GetName() string {
	return "WhiteBit"
}
