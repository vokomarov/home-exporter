package telegram

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vokomarov/home-exporter/config"
)

var Bot *BotService

type BotService struct {
	mu            sync.Mutex
	client        *tgbotapi.BotAPI
	enableWebhook bool
}

func NewBot() (*BotService, error) {
	var err error

	b := BotService{}

	b.client, err = tgbotapi.NewBotAPI(config.Global.TelegramBotToken)
	if err != nil {
		return nil, fmt.Errorf("creating telegram bot instance: %w", err)
	}

	b.client.Debug = false
	b.enableWebhook = false

	log.Printf("Telegram Bot: authorized on account %s", b.client.Self.UserName)

	return &b, nil
}

func (bot *BotService) Listen() error {
	if !bot.enableWebhook {
		return nil
	}

	if err := bot.setWebhook(""); err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	bot.webhookListen(func(update tgbotapi.Update) {
		log.Printf("%+v\n", update)
	})

	return nil
}

func (bot *BotService) Send(message string, chatId int64) error {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	if _, err := bot.client.Send(tgbotapi.MessageConfig{
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

func (bot *BotService) webhookListen(handler func(update tgbotapi.Update)) {
	updates := bot.client.ListenForWebhook("/" + bot.client.Token)

	go http.ListenAndServe("0.0.0.0:80", nil)

	for update := range updates {
		handler(update)
	}
}

func (bot *BotService) setWebhook(url string) error {
	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/%s", url, bot.client.Token))
	if err != nil {
		return fmt.Errorf("creating webhook instance: %w", err)
	}

	if _, err := bot.client.Request(wh); err != nil {
		return fmt.Errorf("making API request to set webhook")
	}

	var info tgbotapi.WebhookInfo

	if info, err = bot.client.GetWebhookInfo(); err != nil {
		return fmt.Errorf("reading webhook info: %w", err)
	}

	if info.LastErrorDate != 0 {
		return fmt.Errorf("error on setting webhook: %s", info.LastErrorMessage)
	}

	log.Printf("Telegram Bot: webhook has been set")

	return nil
}
