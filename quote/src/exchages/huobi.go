package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Huobi struct {
}

type HuobiResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"err-msg"`
	Tick    HuobyTick `json:"tick"`
}

type HuobyTick struct {
	Data []HuobiData
}

type HuobiData struct {
	Price float64 `json:"price"`
}

func (huobi Huobi) Save(track string, base string) {
	var response HuobiResponse
	symbol := strings.ToLower(track + base)
	client := rootConfig.Client

	responseRaw, err := client.Get("https://api.huobi.pro/market/trade?symbol=" + symbol)
	if err != nil {
		log.Println("Huobi connection error. Message: " + err.Error())
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.Status != "ok" {
		log.Println("Huobi request error. Message: " + response.Message)
		return
	}

	quote := model.Quote{
		Exchange:         "Huobi",
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            response.Tick.Data[0].Price,
	}

	db := rootConfig.Database
	db.Save(&quote)
}
