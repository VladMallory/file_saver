package patharchive

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type PathProvider struct {
	log *zap.Logger
}

func NewPathProvider(log *zap.Logger) PathProvider {
	return PathProvider{
		log: log,
	}
}

func closeHelper(closer io.Closer, err *error) {
	if cerr := closer.Close(); cerr != nil {
		*err = errors.Join(*err, cerr)
	}
}

// GetPath получает содержимое файла и отдает каждую строку в slice.
func (p PathProvider) GetPath() (result []string, err error) {
	file, err := os.Open("path.txt")
	if err != nil {
		return nil, err
	}
	defer closeHelper(file, &err)

	scanner := bufio.NewScanner(file)

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

	p.log.Info(
		"файл успешно прочитан",
		zap.Int("lines", len(result)),
	)

	return result, nil
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
