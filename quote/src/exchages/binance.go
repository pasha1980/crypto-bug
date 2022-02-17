package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type Binance struct {
}

type BinancePriceResponse struct {
	Symbol  string `json:"symbol"`
	Price   string `json:"price"`
	Message string `json:"msg"`
}

func (binance Binance) Save(track string, base string) {
	var response BinancePriceResponse

	symbol := track + base
	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + symbol)
	if err != nil {
		log.Println("Binance connection error. Message: " + err.Error())
		return
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if responseRaw.StatusCode != 200 {
		log.Println("Binance return error message: " + response.Message)
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Price, 64)

	quote := model.Quote{
		Exchange:         "Binance",
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            priceFloat,
	}

	db := rootConfig.Database
	db.Save(&quote)
}
