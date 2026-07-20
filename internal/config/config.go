package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config конфигурация.
// nolint: golines
type Config struct {
	TelegramToken  string `env:"TELEGRAM_TOKEN, required"`
	TelegramChatID string `env:"TELEGRAM_CHAT_ID, required"`
	Env            string `env:"ENV"                        env-default:"local"`
	LogLevel       string `env:"LOG_LEVEL"                  env-default:"info"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
