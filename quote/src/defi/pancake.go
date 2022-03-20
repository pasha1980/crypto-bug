package defi

import (
	rootConfig "crypto-bug/config"
	"crypto-bug/model"
	"crypto-bug/service/telegram"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PancakeSwap struct {
}

type PancakeTokenAddrResponse struct {
	Data PancakeTokenAddrResponseData `json:"data"`
}

type PancakeTokenAddrResponseData struct {
	Symbol   string `json:"symbol"`
	PriceBNB string `json:"price_bnb"`
}

func (pancake PancakeSwap) Save(track string, base string) {
	var err error
	track = strings.ToUpper(track)
	base = strings.ToUpper(base)

	trackAddr, exist := rootConfig.Cache.Get("addr.bep20." + track)
	if !exist {
		return
	}

	baseAddr, exist := rootConfig.Cache.Get("addr.bep20." + base)
	if !exist {
		return
	}

	trackBnbPrice, err := pancake.GetBnbPrice(fmt.Sprintf("%v", trackAddr))
	if err != nil {
		telegram.Log("Pancake error: Message: "+err.Error(), "exchange")
		return
	}

	baseBnbPriceInterface, exist := rootConfig.Cache.Get("pancake.token-bnb.price." + base)
	if !exist {
		baseBnbPriceInterface, err = pancake.GetBnbPrice(fmt.Sprintf("%v", baseAddr))
		if err != nil {
			return
		}
	}
	baseBnbPrice, _ := strconv.ParseFloat(fmt.Sprintf("%v", baseBnbPriceInterface), 64)
	if baseBnbPrice == 0.0 {
		telegram.Log("Can't convert pancake var to float64. Token: "+base+fmt.Sprintf(", Price: %v", baseBnbPriceInterface), "exchange")
		return
	}

	newQuote := model.Quote{
		Exchange:      pancake.GetName(),
		Date:          time.Now(),
		BaseCurrency:  base,
		TrackCurrency: track,
		Value:         (trackBnbPrice / baseBnbPrice),
	}
	rootConfig.Database.Save(&newQuote)
}

func (pancake PancakeSwap) GetBnbPrice(addr string) (float64, error) {
	var response PancakeTokenAddrResponse
	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.pancakeswap.info/api/v2/tokens/" + addr)
	if err != nil {
		return 0.0, err
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)
	if responseRaw.StatusCode == 400 {
		return 0.0, errors.New("Pancake request error")
	}

	priceFloat, _ := strconv.ParseFloat(response.Data.PriceBNB, 64)
	rootConfig.Cache.SetTemporary("pancake.token-bnb.price."+strings.ToUpper(response.Data.Symbol), priceFloat)

	return priceFloat, nil
}

func (pancake PancakeSwap) GetName() string {
	return "PancakeSwap"
}
