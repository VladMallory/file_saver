package installer

import (
	"os"
	"path/filepath"
)

// GetAppPath возвращает путь к папке приложения
func GetAppPath() (string, error) {
	// Получаем путь к исполняемому файлу
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Получаем директорию приложения
	return filepath.Dir(execPath), nil
}

// GetPathFilePath возвращает полный путь к файлу path.txt
func GetPathFilePath() (string, error) {
	appDir, err := GetAppPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "path.txt"), nil
}

// ReadPathFile читает пути из файла path.txt в папке приложения
func ReadPathFile() ([]string, error) {
	provider := NewAppPathProvider()
	return provider.GetPath()
}

// WritePathFile записывает пути в файл path.txt в папке приложения
func WritePathFile(paths []string) error {
	return writePathsFile(paths)
}

// PathFileExists проверяет существует ли файл path.txt в папке приложения
func PathFileExists() bool {
	pathFilePath, err := GetPathFilePath()
	if err != nil {
		return false
	}

	_, err = os.Stat(pathFilePath)
	return err == nil
}