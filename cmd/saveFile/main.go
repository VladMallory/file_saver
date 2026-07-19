package main

import (
	"log"

	"saveFile/internal/archive/adapter/outbound/archiveformat/sevenzip"
	patharchive "saveFile/internal/archive/adapter/outbound/path"
	archivecsv "saveFile/internal/archive/service"
	"saveFile/internal/config"
	"saveFile/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	app, err := new()
	if err != nil {
		log.Fatal(err)
	}

	// Сбрасываем буфер логов перед самым выходом из программы.
	// Так как логгер теперь живет внутри app, вызываем Sync через поле структуры.
	defer func() {
		if app.logger != nil {
			_ = app.logger.Sync()
		}
	}()

	app.run()
}

type app struct {
	logger        *zap.Logger
	archiveClient archivecsv.ArchiveService
}

func new() (app, error) {
	cfg, err := config.Load()
	if err != nil {
		return app{}, err
	}

	// ===loger===
	log, err := logger.New(logger.Config{
		ServiceName: "save-file-service",
		Env:         cfg.Env,
		LogLevel:    cfg.LogLevel,
	})
	if err != nil {
		return app{}, err
	}

	log.Info("приложение работает")

	// ===архивация===
	pathArchiveClient := patharchive.NewPathProvider(log)

	sevenzipClient := sevenzip.NewArchiver(log)

	archiveClient := archivecsv.NewArchiveService(log, pathArchiveClient, sevenzipClient)

	//===доставка===

	return app{
		logger:        log,
		archiveClient: archiveClient,
	}, nil
}

func (a app) run() {
	err := a.archiveClient.Run()
	if err != nil {
		return
	}
}
