package telegram

import (
	"log"
	"os"
)

func Log(message string, logType string) {
	tgMessage := TelegramMessage{
		Text:   "[" + logType + "] " + message,
		ChatId: os.Getenv("TELEGRAM_LOG_CHAT_ID"),
	}
	err := tgMessage.Send("log")
	if err != nil {
		ExtraLog(err.Error())
	}

	if logType == "fatal" {
		os.Exit(1)
	}
}

func ExtraLog(message string) {
	f, err := os.OpenFile("/extra.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return // todo
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}
