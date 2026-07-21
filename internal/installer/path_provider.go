package installer

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PathProvider интерфейс для чтения путей из файла
type PathProvider interface {
	GetPath() ([]string, error)
}

// AppPathProvider реализация PathProvider, которая ищет path.txt в папке приложения
type AppPathProvider struct{}

// NewAppPathProvider создает новый AppPathProvider
func NewAppPathProvider() *AppPathProvider {
	return &AppPathProvider{}
}

// GetPath получает содержимое файла path.txt из папки приложения и отдает каждую строку в slice
func (p *AppPathProvider) GetPath() ([]string, error) {
	// Получаем путь к файлу path.txt
	pathFilePath, err := GetPathFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(pathFilePath)
	if err != nil {
		return nil, err
	}
	defer closeHelper(file, &err)

	scanner := bufio.NewScanner(file)

	var result []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		line = expandHome(line)

		// Обрабатывает пути с *
		// Если пришел /root/*. Вернет все что там лежит
		matches, err := filepath.Glob(line)
		if err != nil {
			return nil, err
		}

		if len(matches) == 0 {
			result = append(result, line)
		} else {
			result = append(result, matches...)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func closeHelper(closer io.Closer, err *error) {
	if cerr := closer.Close(); cerr != nil {
		*err = errors.Join(*err, cerr)
	}
}

func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return strings.Replace(path, "~", home, 1)
}