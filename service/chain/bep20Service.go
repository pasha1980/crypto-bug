package chain

import (
	rootConfig "crypto-bug/config"
	quoteConfig "crypto-bug/quoteConf"
	"crypto-bug/service/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type Bep20 struct {
}

type PancakeTokenResponse struct {
	Data map[string]PancakeTokenData `json:"data"`
}

type PancakeTokenData struct {
	Symbol string `json:"symbol"`
}

// todo: Найти более полную бд токенов bep20
func (bep20 Bep20) GetTokens() error {
	var response PancakeTokenResponse

	client := rootConfig.Client
	responseRaw, err := client.Get("https://api.pancakeswap.info/api/v2/tokens")
	if err != nil {
		return err
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)

	var (
		isset bool
		addr  interface{}
	)

	for _, trackToken := range quoteConfig.CurrenciesToTrack {
		trackToken = strings.ToUpper(trackToken)
		isset, addr = utils.InArray(PancakeTokenData{Symbol: trackToken}, response.Data)
		if !isset {
			continue
		}
		addr = fmt.Sprintf("%v", addr)
		rootConfig.Cache.Set("addr.bep20."+trackToken, addr)
	}

	for _, baseToken := range quoteConfig.BaseCurrencies {
		baseToken = strings.ToUpper(baseToken)
		isset, addr = utils.InArray(PancakeTokenData{Symbol: baseToken}, response.Data)
		if !isset {
			continue
		}
		addr = fmt.Sprintf("%v", addr)
		rootConfig.Cache.Set("addr.bep20."+baseToken, addr)
	}

	return nil
}
