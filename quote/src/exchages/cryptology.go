package exchages

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/quote"
	"crypto-bug/service/telegram"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Cryptology struct {
}

type CryptologyResponse struct {
	Status string                    `json:"status"`
	Error  CryptologyResponseError   `json:"error"`
	Data   [1]CryptologyResponseData `json:"data"`
}

type CryptologyResponseError struct {
	Message string `json:"message"`
}

type CryptologyResponseData struct {
	Price string `json:"price"`
}

const cryptologyNeedExceptionMessage = "invalid trade_pair"
const cryptologyNeedToSleepMessage = "Too many requests"

func (cryptology Cryptology) Save(track string, base string) {
Start:

	var response CryptologyResponse
	symbol := strings.ToUpper(track + "_" + base)
	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.cryptology.com/v1/public/get-trades?limit=1&trade_pair=" + symbol)
	if err != nil {
		telegram.Log("Cryptology connection error. Message: "+err.Error(), "exchange")
		return
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	if response.Status != "OK" {
		if response.Error.Message == cryptologyNeedExceptionMessage {
			quote.ProcessException(cryptology.GetName(), track, base)
		} else if response.Error.Message == cryptologyNeedToSleepMessage {
			time.Sleep(2 * time.Second)
			goto Start
		} else {
			telegram.Log("Cryptology return error message: "+response.Error.Message, "exchange")
		}
		return
	}

	priceFloat, _ := strconv.ParseFloat(response.Data[0].Price, 64)

	newQuote := model.Quote{
		Exchange:      cryptology.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         priceFloat,
	}

	db := rootConfig.Database
	db.Save(&newQuote)
}

func (cryptology Cryptology) GetName() string {
	return "Cryptology"
}
