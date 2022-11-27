package telegram

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vokomarov/home-exporter/config"
)

var bot *tgbotapi.BotAPI
var enableWebhook = false

func Listen() error {
	var err error

	bot, err = tgbotapi.NewBotAPI(config.Global.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("creating telegram bot instance: %w", err)
	}

	log.Printf("Telegram Bot: authorized on account %s", bot.Self.UserName)

	bot.Debug = true

	if enableWebhook {
		if err := setWebhook(""); err != nil {
			return fmt.Errorf("set webhook: %w", err)
		}

		webhookListen(func(update tgbotapi.Update) {
			log.Printf("%+v\n", update)
		})
	}

	return nil
}

func Send(message string, chatId int64) error {
	if _, err := bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatId,
		},
		Text:      message,
		ParseMode: "HTML",
	}); err != nil {
		return fmt.Errorf("telegram bot send messaage failed: %v", err)
	}

	return nil
}

func webhookListen(handler func(update tgbotapi.Update)) {
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:80", nil)

	for update := range updates {
		handler(update)
	}
}

func setWebhook(url string) error {
	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/%s", url, bot.Token))
	if err != nil {
		return fmt.Errorf("creating webhook instance: %w", err)
	}

	if _, err := bot.Request(wh); err != nil {
		return fmt.Errorf("making API request to set webhook")
	}

	var info tgbotapi.WebhookInfo

	if info, err = bot.GetWebhookInfo(); err != nil {
		return fmt.Errorf("reading webhook info: %w", err)
	}

	if info.LastErrorDate != 0 {
		return fmt.Errorf("error on setting webhook: %s", info.LastErrorMessage)
	}

	log.Printf("Telegram Bot: webhook has been set")

	return nil
}
