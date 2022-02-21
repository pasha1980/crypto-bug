package parserService

import (
	"bytes"
	rootConfig "crypto-bug/config"
	"encoding/json"
	"errors"
	"os"
)

type TelegramMessage struct {
	Text   string `json:"text"`
	ChatId string `json:"chat_id"`
}

var BotTypeToken = map[string]string{
	"log":  "TELEGRAM_LOG_BOT_TOKEN",
	"main": "TELEGRAM_BOT_TOKEN",
}

type TelegramResponse struct {
	Description string `json:"description"`
}

func (message TelegramMessage) Send(botType string) error {
	var response TelegramResponse

	TokenEnvName := BotTypeToken[botType]
	client := rootConfig.Client
	token := os.Getenv(TokenEnvName)

	body, _ := json.Marshal(message)
	responseRaw, err := client.Post(
		"https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		ExtraLog("Sending message error. Message: " + err.Error())
		return err
	}

	defer responseRaw.Body.Close()
	_ = json.NewDecoder(responseRaw.Body).Decode(&response)
	if responseRaw.StatusCode != 200 {
		var fullResponse string
		_ = json.NewDecoder(responseRaw.Body).Decode(&fullResponse)
		ExtraLog("Telegram return non 200 response. Message: " + response.Description + ". Response: " + fullResponse)
		return errors.New(response.Description)
	}
	return nil
}

func (message TelegramMessage) SendToMe() {
	myChatId := os.Getenv("TELEGRAM_MY_CHAT_ID")
	message.ChatId = myChatId
	_ = message.Send("main")
}
