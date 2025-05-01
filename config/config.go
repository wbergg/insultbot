package config

import (
	"encoding/json"
	"os"
)

type TelegramConfig struct {
	Enabled   bool   `json:"enabled,string"`
	TgAPIKey  string `json:"tgAPIkey"`
	TgChannel string `json:"tgChannel"`
}

type IRCConfig struct {
	Enabled  bool   `json:"enabled,string"`
	Server   string `json:"server"`
	Nick     string `json:"nick"`
	User     string `json:"user"`
	Channel  string `json:"channel"`
	Password string `json:"password"`
}

type Config struct {
	Telegram TelegramConfig `json:"Telegram"`
	IRC      IRCConfig      `json:"IRC"`
}

var Loaded Config

func LoadConfig(filepath string) (Config, error) {
	var c Config

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	Loaded = c

	return c, nil
}
