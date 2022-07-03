package quoteConf

import (
	"crypto-bug/quote/src/exchages"
)

type Exchange interface {
	Save(track string, base string)
	GetName() string
}

var Exchanges = []Exchange{
	exchages.Binance{},
	exchages.WhiteBit{},
	//exchages.Coinlist{},
	exchages.ByBit{},
	exchages.Huobi{},
	//exchages.Cryptology{},
	//defi.PancakeSwap{},
}

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
	"1INCH",
	//"LUNA",
	"ATOM",
	"LINK",
	"BCH",
	"SHIB",
	"FTT",
	"CRO",
	"NEAR",
	"ALGO",
	"BCH",
	"XMR",
}

var BaseCurrencies = []string{
	"USDT",
	"USDC",
	"BUSD",
	//"TUSD",
	//"UST",
}
