package installer

import (
	"os"
	"path/filepath"
	"strings"
)

// writePathsFile записывает (перезаписывает) path.txt в текущую рабочую директорию.
// Каждый путь — на отдельной строке.
// Формат совместим с PathProvider из archive/adapter/outbound/path.
func writePathsFile(paths []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dst := filepath.Join(wd, "path.txt")
	content := strings.Join(paths, "\n") + "\n"

	return os.WriteFile(dst, []byte(content), 0o644)
}
