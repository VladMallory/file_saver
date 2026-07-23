package installer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

func Run() error {
	token, chatID, err := collectTelegramConfig()
	if err != nil {
		return err
	}

	paths, err := collectPaths()
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	appDir := filepath.Dir(exePath)

	if err := saveFiles(appDir, token, chatID, paths); err != nil {
		return err
	}

	var addCron bool

	err = huh.NewConfirm().
		Title("Добавить задачу в cron").
		Description("Запуск каждый день в 04:00").
		Value(&addCron).
		Run()
	if err != nil {
		return err
	}

	if addCron {
		if err := setupCron(exePath); err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintln(os.Stdout, "\n🎉 Настройка успешно завершена!")
	_, _ = fmt.Fprintf(
		os.Stdout,
		"Конфигурационные файлы сохранены рядом с бинарником: %s\n",
		appDir,
	)

	return nil
}

const da = "Берется из @BotFather. Нужно будет создать бота, или если уже есть, " +
	"взять токен из существующего -> /mybots -> имя бота -> API TOKEN -> скопировать и вставить сюда"

func collectTelegramConfig() (string, string, error) {
	var token, chatID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Телеграм бот токен").
				Description(da).
				Value(&token),

			huh.NewInput().
				Title("Telegram Chat ID").
				Description("Берется из @Getmyid_bot. Вроде бы не требует подписку и не показывает рекламу").
				Value(&chatID),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", err
	}

	return token, chatID, nil
}

func collectPaths() ([]string, error) {
	var paths []string

	for {
		var inputPath string

		pathInput := huh.NewInput().
			Title(fmt.Sprintf("Добавление пути №%d", len(paths)+1)).
			Description("Введите stop чтобы закончить с путями которые будете бекапить").
			Value(&inputPath)

		if err := pathInput.Run(); err != nil {
			return nil, err
		}

		trimmed := strings.TrimSpace(inputPath)

		if strings.EqualFold(trimmed, "stop") {
			if len(paths) == 0 {
				return nil, errors.New("укажите хотя бы один путь")
			}

			break
		}

		if trimmed != "" {
			paths = append(paths, trimmed)
		}
	}

	return paths, nil
}

func saveFiles(appDir, token, chatID string, paths []string) error {
	envPath := filepath.Join(appDir, ".env")
	envContent := fmt.Sprintf("TELEGRAM_TOKEN=%s\nTELEGRAM_CHAT_ID=%s\n", token, chatID)

	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		return fmt.Errorf("ошибка сохранения .env: %w", err)
	}

	pathFilePath := filepath.Join(appDir, "path.txt")
	pathContent := strings.Join(paths, "\n")
	if len(paths) > 0 {
		pathContent += "\n"
	}

	if err := os.WriteFile(pathFilePath, []byte(pathContent), 0o600); err != nil {
		return fmt.Errorf("ошибка сохранения %s: %w", pathFilePath, err)
	}

	return nil
}

func setupCron(exePath string) error {
	cronLine := fmt.Sprintf("0 4 * * * %s run\n", exePath)
	cmd := exec.Command(
		"sh", "-c",
		fmt.Sprintf("(crontab -l 2>/dev/null; echo '%s') | crontab -",
			cronLine),
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ошибка добавления cron: %w\n%s", err, out)
	}

	return nil
}
