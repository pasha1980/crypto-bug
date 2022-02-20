package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
	"encoding/json"
	"strconv"
	"time"
)

type Binance struct {
}

type BinanceResponse struct {
	Symbol  string `json:"symbol"`
	Price   string `json:"price"`
	Message string `json:"msg"`
}

func (binance Binance) Save(track string, base string) {
	var response BinanceResponse

	symbol := track + base
	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + symbol)
	if err != nil {
		service.Log("Binance connection error. Message: "+err.Error(), "exchange")
		return
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if responseRaw.StatusCode != 200 {
		service.Log("Binance return error message: "+response.Message, "exchange")
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Price, 64)

	quote := model.Quote{
		Exchange:      binance.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         priceFloat,
	}

	db := rootConfig.Database
	db.Save(&quote)
}

func (binance Binance) GetName() string {
	return "Binance"
}
