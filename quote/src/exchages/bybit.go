package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
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

func (byBit ByBit) Save(track string, base string) {
	var response ByBitResponse
	symbol := track + base
	client := rootConfig.Client

	responseRaw, err := client.Get("https://api.bybit.com/v2/public/tickers?symbol=" + symbol)
	if err != nil {
		service.Log("Bybit connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.ReturnCode != 0 {
		service.Log("Bybit request error. Message: "+response.ReturnMessage, "exchange")
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Result[0].Price, 64)

	quote := model.Quote{
		Exchange:         byBit.GetName(),
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            priceFloat,
	}

	db := rootConfig.Database
	db.Save(&quote)
}

func (bybit ByBit) GetName() string {
	return "ByBit"
}
