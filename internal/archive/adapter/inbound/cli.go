package archiveinboud

import (
	"flag"
	"time"

	"go.uber.org/zap"
)

type archiveUseCase interface {
	Run(savePathFile string) error
}

type Handler struct {
	log            *zap.Logger
	archiveUseCase archiveUseCase
}

func NewHandler(log *zap.Logger, uc archiveUseCase) Handler {
	return Handler{
		log:            log,
		archiveUseCase: uc,
	}
}

func (h Handler) Execute(args []string, savePathFile string) error {
	start := time.Now()

	defer func() {
		duration := time.Since(start)
		h.log.Info(
			"Архивация успешно выполнена",
			zap.Duration("duration", duration),
		)
	}()

	fs := flag.NewFlagSet("saveFile", flag.ContinueOnError)
	force := fs.Bool("force", false, "принудительный запуск")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := h.archiveUseCase.Run(savePathFile); err != nil {
		h.log.Error("Ошибка при выполнении архивации", zap.Error(err))
		return err
	}

	h.log.Info(
		"выполнена cli команда",
		zap.Bool("forece", *force),
	)

	return nil
}
