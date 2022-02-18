package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"encoding/json"
	"log"
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
		log.Println("Bybit connection error. Message: " + err.Error())
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.ReturnCode != 0 {
		log.Println("Bybit request error. Message: " + response.ReturnMessage)
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
