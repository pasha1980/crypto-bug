package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/quote"
	"crypto-bug/service/telegram"
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

const huobiReturnMessageNeedException = "invalid symbol"

func (huobi Huobi) Save(track string, base string) {
	var response HuobiResponse
	symbol := strings.ToLower(track + base)
	client := rootConfig.Client

	responseRaw, err := client.Get("https://api.huobi.pro/market/trade?symbol=" + symbol)
	if err != nil {
		telegram.Log("Huobi connection error. Message: "+err.Error(), "exchange")
		return
	}
	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.Status != "ok" {
		if response.Message == huobiReturnMessageNeedException {
			quote.ProcessException(huobi.GetName(), track, base)
		} else {
			telegram.Log("Huobi request error. Message: "+response.Message, "exchange")
		}
		return
	}

	newQuote := model.Quote{
		Exchange:      huobi.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         response.Tick.Data[0].Price,
	}

	db := rootConfig.Database
	db.Save(&newQuote)
}

func (huobi Huobi) GetName() string {
	return "Huobi"
}
