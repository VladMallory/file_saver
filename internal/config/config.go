package config

import (
	"os"
	"path/filepath"

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
	// 1. Сначала пробуем взять .env из текущей папки (для локальной разработки)
	_ = godotenv.Load(".env")

	// 2. Если не получилось, пробуем взять .env рядом с бинарником (для продакшена)
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		envPath := filepath.Join(exeDir, ".env")
		_ = godotenv.Load(envPath)
	}

	// 3. Читаем переменные (из системы или из .env)
	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
