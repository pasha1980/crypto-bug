package main

import (
	"crypto-bug/config"
	"crypto-bug/parser"
	"crypto-bug/quote"
	"time"
)

func main()  {
	config.Initialization()
	for range time.Tick(time.Minute) {
		go quote.SaveQuotes()
		go parser.Init()
	}
}
