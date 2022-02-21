package main

import (
	"crypto-bug/config"
	"crypto-bug/migrations"
	"crypto-bug/parser"
	"crypto-bug/parser/src/service"
	"crypto-bug/quote"
	"os"
	"time"
)

func main() {
	config.Initialization()
	parserService.Log("Application started", "info")
	migrations.Migrate()
	repeat, err := time.ParseDuration(os.Getenv("REPEAT_TIME"))
	if err != nil {
		parserService.Log(err.Error(), "fatal")
	}

	if repeat == 0 {
		repeat = time.Minute
	}

	for range time.Tick(repeat) {
		go quote.Init()
		go parser.Init()
	}
}
