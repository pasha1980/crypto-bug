package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type WhiteBit struct {
}

type WhiteBitResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Result  []WhiteBitMarketResult `json:"result"`
}

type WhiteBitMarketResult struct {
	Price string `json:"price"`
}

func (whiteBit WhiteBit) Save(track string, base string) {
	var response WhiteBitResponse
	market := track + "_" + base
	client := rootConfig.Client

	responseRaw, err := client.Get("https://whitebit.com/api/v1/public/history?lastId=1&limit=1&market=" + market)
	if err != nil {
		log.Println("WhiteBit connection error. Message: " + err.Error())
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if !response.Success {
		log.Println("WhiteBit request error. Message: " + response.Message)
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Result[0].Price, 64)

	quote := model.Quote{
		Exchange:         "WhiteBit",
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            priceFloat,
	}

	db := rootConfig.Database
	db.Save(&quote)
}