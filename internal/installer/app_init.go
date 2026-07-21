package installer

import (
	"os"
)

// EnsurePathFileExists проверяет существует ли файл path.txt в папке приложения, и если нет - создает его.
func EnsurePathFileExists() error {
	// Получаем путь к файлу path.txt
	pathFilePath, err := GetPathFilePath()
	if err != nil {
		return err
	}

	// Проверяем существует ли файл
	if _, err := os.Stat(pathFilePath); os.IsNotExist(err) {
		// Создаем пустой файл
		return os.WriteFile(pathFilePath, []byte(""), 0o644)
	}

	return nil
}

// GetCurrentPathFile возвращает содержимое файла path.txt из папки приложения.
func GetCurrentPathFile() ([]string, error) {
	provider := NewAppPathProvider()

	return provider.GetPath()
}
