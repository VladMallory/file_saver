package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeEnvFile записывает (перезаписывает) .env файл рядом с исполняемым бинарём.
// Это соответствует логике config.Load() — в production она ищет .env рядом с бинарём.
func writeEnvFile(token, chatID string) error {
	exePath, err := os.Executable()
	if err != nil {
		exePath = ".env"
	} else {
		exePath = filepath.Join(filepath.Dir(exePath), ".env")
	}

	content := fmt.Sprintf("TELEGRAM_TOKEN=%s\nTELEGRAM_CHAT_ID=%s\n", token, chatID)

	return os.WriteFile(exePath, []byte(content), 0o644)
}
