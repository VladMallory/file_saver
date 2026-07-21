package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// envPath возвращает путь к .env файлу рядом с исполняемым бинарём.
func envPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return ".env"
	}
	return filepath.Join(filepath.Dir(exePath), ".env")
}

// readEnvValue читает из .env файла значение переменной по ключу.
// Возвращает пустую строку, если файла нет или ключ не найден.
func readEnvValue(key string) string {
	data, err := os.ReadFile(envPath())
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

// writeEnvFile записывает (перезаписывает) .env файл рядом с исполняемым бинарём.
// Это соответствует логике config.Load() — в production она ищет .env рядом с бинарём.
func writeEnvFile(token, chatID string) error {
	content := fmt.Sprintf("TELEGRAM_TOKEN=%s\nTELEGRAM_CHAT_ID=%s\n", token, chatID)

	err := os.WriteFile(envPath(), []byte(content), 0o644)
	if err != nil {
		return err
	}

	return nil
}
