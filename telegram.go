package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Telegram struct {
	BotKey string
	ChatID string
}

func NewTelegram(BotID, ChatID string) *Telegram {

	return &Telegram{
		BotKey: BotID,
		ChatID: ChatID,
	}
}
func (t *Telegram) SendTelegramMessage(message string) {
	u, _ := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s", t.BotKey, t.ChatID))
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("text", message)
	u.RawQuery = q.Encode()
	response, err := http.Get(u.String())
	if err != nil {
		log.Printf("Cannot send telegram message %v", err)
		return
	}
	buf, _ := ioutil.ReadAll(response.Body)
	log.Print(string(buf))
}
