package sevenzip

import (
	"errors"
	"fmt"
	"os/exec"

	archivecore "saveFile/internal/archive/domain"

	"go.uber.org/zap"
)

type Archiver struct {
	bin string
	log *zap.Logger
}

func NewArchiver(log *zap.Logger) Archiver {
	return Archiver{
		bin: "7z",
		log: log,
	}
}

func (a Archiver) Run(paths []string, saveFilePath string) (string, error) {
	if len(paths) == 0 {
		return "", archivecore.ErrNoPaths
	}

	args := []string{
		"a",          // Создать архив (Add)
		"-t7z",       // Использовать формат 7z
		"-m0=lzma2",  // Алгоритм сжатия LZMA2
		"-mx=9",      // Максимальный уровень сжатия (0-9)
		"-mmt=on",    // Использовать все доступные ядра процессора
		"-ms=on",     // Solid-архив  лучшее сжатие для набора похожих файлов
		"-md=256m",   // Размер словаря LZMA2 (256 МБ)
		"-mfb=273",   // Максимальное количество Fast Bytes (улучшает сжатие)
		saveFilePath, // Путь к создаваемому архиву
	}

	args = append(args, paths...)

	cmd := exec.Command(a.bin, args...)

	if outBytes, err := cmd.CombinedOutput(); err != nil {
		// 7z возвращает exit code 1 при WARNING (нет файлов)
		// Это не критично — архив создаётся, просто пустой
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			// a.log.Warn(
			// 	"7z завершился с предупреждением",
			// 	zap.ByteString("output", outBytes),
			// )
		} else {
			return "", fmt.Errorf("%w: %s", err, outBytes)
		}
	}

	return saveFilePath, nil
}
