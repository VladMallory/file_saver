package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// installCronJob добавляет задание в crontab текущего пользователя.
// Существующие записи сохраняются; старые записи #saveFile заменяются новыми.
func installCronJob(cfg CronSettings) error {
	schedule := cronSchedule(cfg.Time, cfg.Interval)
	if schedule == "" {
		parts := strings.Split(cfg.Time, ":")
		if len(parts) == 2 {
			schedule = fmt.Sprintf("%s %s * * *", parts[1], parts[0])
		} else {
			schedule = "0 2 * * *" // fallback: ежедневно в 02:00
		}
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("не удалось определить путь к бинарю: %w", err)
	}

	cronLine := fmt.Sprintf("%s %s run #saveFile", schedule, exePath)

	currentCmd := exec.Command("crontab", "-l")
	out, _ := currentCmd.Output()

	existing := strings.TrimSpace(string(out))
	excisting := splitLinesPreserve(existing)

	// Удаляем старые записи #saveFile
	var newLines []string
	for _, line := range excisting {
		if !strings.Contains(line, "#saveFile") {
			newLines = append(newLines, line)
		}
	}

	newLines = append(newLines, cronLine)
	newContent := strings.Join(newLines, "\n") + "\n"

	writeCmd := exec.Command("crontab", "-")
	writeCmd.Stdin = strings.NewReader(newContent)

	if err := writeCmd.Run(); err != nil {
		return fmt.Errorf("не удалось обновить crontab: %w", err)
	}

	return nil
}

// splitLinesPreserve разбивает строку на строки, сохраняя пустые как отдельные элементы.
// Нужна чтобы корректно обработать случай, когда crontab полностью пуст.
func splitLinesPreserve(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// cronSchedule преобразует время и интервал в cron-выражение из 5 полей.
func cronSchedule(time, interval string) string {
	parts := strings.Split(time, ":")
	if len(parts) != 2 {
		return ""
	}
	minute, hour := parts[1], parts[0]

	switch strings.ToLower(interval) {
	case "daily":
		return fmt.Sprintf("%s %s * * *", minute, hour)
	case "weekly":
		return fmt.Sprintf("%s %s * * 0", minute, hour) // воскресенье
	case "monthly":
		return fmt.Sprintf("%s %s 1 * *", minute, hour) // 1-е число
	default:
		return fmt.Sprintf("%s %s * * *", minute, hour) // fallback: daily
	}
}
