package config

import (
	"crypto-bug/quote/src/exchages"
	"time"
)

var HoursToSaveQuotes time.Duration = 1

var CurrenciesToTrack = []string{
	"BTC",
	"ETH",
	"ADA",
	"BNB",
	"MATIC",
}

var BaseCurrencies = []string{
	"USDT",
}

var Exchanges = []exchages.Exchange{
	exchages.Binance{},
	exchages.WhiteBit{},
	exchages.Coinlist{},
}
