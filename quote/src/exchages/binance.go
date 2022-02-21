package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/parserService"
	"crypto-bug/quote/src/service"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Binance struct {
}

type BinanceResponse struct {
	Symbol  string `json:"symbol"`
	Price   string `json:"price"`
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

const binanceReturnCodeNeedException = -1121

func (binance Binance) Save(track string, base string) {
	var response BinanceResponse

	symbol := strings.ToUpper(track + base)
	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + symbol)
	if err != nil {
		parserService.Log("Binance connection error. Message: "+err.Error(), "exchange")
		return
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if responseRaw.StatusCode != 200 {
		if responseRaw.StatusCode == 400 && response.Code == binanceReturnCodeNeedException {
			service.ProcessException(binance.GetName(), track, base)
		} else {
			parserService.Log("Binance return error message: "+response.Message, "exchange")
		}
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Price, 64)

	newQuote := model.Quote{
		Exchange:      binance.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         priceFloat,
	}

	db := rootConfig.Database
	db.Save(&newQuote)
}

func (binance Binance) GetName() string {
	return "Binance"
}
