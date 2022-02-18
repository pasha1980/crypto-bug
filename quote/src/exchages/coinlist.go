package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type Coinlist struct {
}

type CoinlistResponse struct {
	LastTrade CoinistTrade `json:"last_trade"`
	Message   string       `json:"message"`
}

type CoinistTrade struct {
	Price string `json:"price"`
}

func (coinist Coinlist) Save(track string, base string) {
	var response CoinlistResponse
	symbol := track + "-" + base
	client := rootConfig.Client

	responseRaw, err := client.Get("https://trade-api.coinlist.co/v1/symbols/" + symbol + "/summary")
	if err != nil {
		log.Println("Coinlist connection error. Message: " + err.Error())
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if responseRaw.StatusCode != 200 {
		log.Println("Coinlist request error. Message: " + response.Message)
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.LastTrade.Price, 64)

	quote := model.Quote{
		Exchange:         coinist.GetName(),
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            priceFloat,
	}

	db := rootConfig.Database
	db.Save(&quote)
}

func (coinlist Coinlist) GetName() string {
	return "Coinlist"
}
