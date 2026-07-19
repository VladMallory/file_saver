package sevenzip

import (
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

func (a Archiver) Run(paths []string) (string, error) {
	if len(paths) == 0 {
		return "", archivecore.ErrNoPaths
	}

	outPath := "backup.7z"

	args := []string{
		"a",         // Создать архив (Add)
		"-t7z",      // Использовать формат 7z
		"-m0=lzma2", // Алгоритм сжатия LZMA2
		"-mx=9",     // Максимальный уровень сжатия (0-9)
		"-mmt=on",   // Использовать все доступные ядра процессора
		"-ms=on",    // Solid-архив  лучшее сжатие для набора похожих файлов
		"-md=256m",  // Размер словаря LZMA2 (256 МБ)
		"-mfb=273",  // Максимальное количество Fast Bytes (улучшает сжатие)
		outPath,     // Путь к создаваемому архиву
	}

	args = append(args, paths...)

	cmd := exec.Command(a.bin, args...)

	if outBytes, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%w: %s", err, outBytes)
	}

	return outPath, nil
}
