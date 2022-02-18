package service

import (
	"bytes"
	rootConfig "crypto-bug/config"
	"encoding/json"
	"log"
	"os"
)

type TelegramMessage struct {
	Text   string `json:"text"`
	ChatId string `json:"chat_id"`
}

type TelegramResponse struct {
	Description string `json:"description"`
}

func (message TelegramMessage) Send() {
	var response TelegramResponse

	client := rootConfig.Client
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	body, _ := json.Marshal(message)
	responseRaw, err := client.Post(
		"https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		log.Println("Sending message error. Message: " + err.Error())
		return
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)
	if responseRaw.StatusCode != 200 {
		log.Println("Telegram return non 200 response. Message: " + response.Description)
	}
}

func (message TelegramMessage) SendToMe() {
	myChatId := os.Getenv("TELEGRAM_MY_CHAT_ID")
	message.ChatId = myChatId
	message.Send()
}
