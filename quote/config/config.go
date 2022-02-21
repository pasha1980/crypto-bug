package config

import (
	"time"
)

var HoursToSaveQuotes time.Duration = 1

var CurrenciesToTrack = []string{
	"BTC",
	"ETH",
	"ADA",
	"BNB",
	"MATIC",
	"XRP",
	"SOL",
	"DOT",
	"LTC",
	"TRX",
	"UNI",
}

var BaseCurrencies = []string{
	"USDT",
	"USDC",
	"BUSD",
}
