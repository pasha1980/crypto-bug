package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/parserService"
	"crypto-bug/quote/src/service"
	"encoding/json"
	"strconv"
	"time"
)

type ByBit struct {
}

type ByBitResponse struct {
	ReturnCode    int           `json:"ret_code"`
	ReturnMessage string        `json:"ret_msg"`
	Result        []ByBitResult `json:"result"`
}

type ByBitResult struct {
	Price string `json:"last_price"`
}

const bybitReturnCodeNeedException = 10001

func (byBit ByBit) Save(track string, base string) {
	var response ByBitResponse
	symbol := track + base
	client := rootConfig.Client

	responseRaw, err := client.Get("https://api.bybit.com/v2/public/tickers?symbol=" + symbol)
	if err != nil {
		parserService.Log("Bybit connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.ReturnCode != 0 {
		if response.ReturnCode == bybitReturnCodeNeedException {
			service.ProcessException(byBit.GetName(), track, base)
		} else {
			parserService.Log("Bybit request error. Message: "+response.ReturnMessage, "exchange")
		}
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Result[0].Price, 64)

	newQuote := model.Quote{
		Exchange:      byBit.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         priceFloat,
	}

	db := rootConfig.Database
	db.Save(&newQuote)
}

func (bybit ByBit) GetName() string {
	return "ByBit"
}
