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
	config.Initialization()
	telegram.Log("Application started", "info")
	migrations.Migrate()
	repeat, err := time.ParseDuration(os.Getenv("REPEAT_TIME"))
	if err != nil {
		telegram.Log(err.Error(), "fatal")
	}

	if repeat == 0 {
		repeat = time.Minute
	}

	for range time.Tick(repeat) {
		quote.Init()
		parser.Init()
	}
}
