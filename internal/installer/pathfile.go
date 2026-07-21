package installer

import (
	"os"
	"path/filepath"
	"strings"
)

// writePathsFile записывает (перезаписывает) path.txt в папку приложения.
// Каждый путь — на отдельной строке.
// Формат совместим с PathProvider из archive/adapter/outbound/path.
func writePathsFile(paths []string) error {
	// Получаем путь к исполняемому файлу
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Получаем директорию приложения
	appDir := filepath.Dir(execPath)

	dst := filepath.Join(appDir, "path.txt")
	content := strings.Join(paths, "\n") + "\n"

	return os.WriteFile(dst, []byte(content), 0o644)
}
