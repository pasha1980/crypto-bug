package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
	"encoding/json"
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
		service.Log("Coinlist connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if responseRaw.StatusCode != 200 {
		service.Log("Coinlist request error. Message: "+response.Message, "exchange")
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
