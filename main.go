package main

import (
	"crypto-bug/config"
	"crypto-bug/migrations"
	"crypto-bug/parser"
	"crypto-bug/quote"
	"log"
	"os"
	"time"
)

func main() {
	config.Initialization()
	migrations.Migrate()
	repeat, err := time.ParseDuration(os.Getenv("REPEAT_TIME"))
	if err != nil {
		log.Fatal(err)
	}

	if repeat == 0 {
		repeat = time.Minute
	}

	for range time.Tick(repeat) {
		go quote.Init()
		go parser.Init()
	}
}
