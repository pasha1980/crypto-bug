package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/parser/src/service"
	"encoding/json"
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
		service.Log("Huobi connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.Status != "ok" {
		service.Log("Huobi request error. Message: "+response.Message, "exchange")
		return
	}

	quote := model.Quote{
		Exchange:         huobi.GetName(),
		Date:             time.Now(),
		BaseCurrency:     base,
		TrackingCurrency: track,
		Value:            response.Tick.Data[0].Price,
	}

	db := rootConfig.Database
	db.Save(&quote)
}

func (huobi Huobi) GetName() string {
	return "Huobi"
}
