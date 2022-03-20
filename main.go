package main

import (
	"crypto-bug/config"
	"crypto-bug/migrations"
	"crypto-bug/parser"
	"crypto-bug/quote"
	"crypto-bug/service/telegram"
	"os"
	"time"
)

func main() {
	var err error
	config.Initialization()
	telegram.Log("Application started", "info")
	migrations.Migrate()
	repeat := os.Getenv("REPEAT_TIME")

	var repeatDuration = time.Minute
	if repeat != "" {
		repeatDuration, err = time.ParseDuration(repeat)
		if err != nil {
			telegram.Log(err.Error(), "fatal")
		}
	}

	for range time.Tick(repeatDuration) {
		quote.Init()
		parser.Init()
		config.Cache.Clear()
	}
}
