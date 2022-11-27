package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const DefaultConfigPathEnvName = "CONFIG_PATH"
const TelegramBotTokenEnvName = "TG_BOT_TOKEN"
const DefaultConfigPath = "config.yml"

var Global Config

type Config struct {
	TelegramBotToken string `yaml:"telegramBotToken"`
	Homes            []Home `yaml:"homes"`
}

type Home struct {
	Name           string         `yaml:"name"`
	TelegramChatId int64          `yaml:"telegramChatId"`
	InternetStatus InternetStatus `yaml:"internetStatus"`
}

type InternetStatus struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Method   string `yaml:"method"`
	Retries  int    `yaml:"retries"`
	Timeout  int    `yaml:"timeout"`  // in seconds
	Interval int    `yaml:"interval"` // in seconds
}

func Load() error {
	configPath := os.Getenv(DefaultConfigPathEnvName)
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	Global.Homes = make([]Home, 0)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &Global)
	if err != nil {
		return fmt.Errorf("parsing yaml file: %w", err)
	}

	overrideFromEnv()

	return nil
}

func overrideFromEnv() {
	if telegramBotTokenEnv := os.Getenv(TelegramBotTokenEnvName); telegramBotTokenEnv != "" {
		Global.TelegramBotToken = telegramBotTokenEnv
	}
}
