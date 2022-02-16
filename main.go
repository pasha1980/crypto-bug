package main

import (
	"crypto-bug/config"
	"crypto-bug/migrations"
	"crypto-bug/parser"
	"crypto-bug/quote"
	"time"
)

func main() {
	config.Initialization()
	migrations.Migrate()
	for range time.Tick(time.Minute) {
		go quote.Init()
		go parser.Init()
	}
}
