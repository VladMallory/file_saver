package main

import (
	"context"
	"log"
	"os"
	"saveFile/internal/archive/adapter/outbound/archiveformat/sevenzip"
	"saveFile/internal/config"
	"saveFile/pkg/logger"
	"time"

	archiveinboud "saveFile/internal/archive/adapter/inbound"

	patharchive "saveFile/internal/archive/adapter/outbound/path"
	archiveusecase "saveFile/internal/archive/service"

	deliverytelegram "saveFile/internal/deliveryArchive/adapter"
	delivery "saveFile/internal/deliveryArchive/domain"
	deliveryusecase "saveFile/internal/deliveryArchive/service"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run — собирает и запускает приложение.
func run() error {
	app, err := new()
	if err != nil {
		return err
	}
	defer func() {
		_ = app.log.Sync()
	}()

	return app.start()
}

type app struct {
	log            *zap.Logger
	cliHandler     archiveinboud.Handler
	deliveryClient deliveryusecase.DeliveryService
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

	// ===архивация===
	pathArchiveClient := patharchive.NewPathProvider(log)
	sevenzipClient := sevenzip.NewArchiver(log)
	archiveClient := archiveusecase.NewArchiveService(log, pathArchiveClient, sevenzipClient)
	cliHandler := archiveinboud.NewHandler(log, archiveClient)

	// ===доставка===
	telegramClient, err := deliverytelegram.NewTelegramSender(
		log, cfg.TelegramToken,
		cfg.TelegramChatID,
	)
	if err != nil {
		return app{}, err
	}

	deliveryClient := deliveryusecase.NewDeliveryService(log, telegramClient)

	log.Info("приложение собралось")

	return app{
		log:            log,
		cliHandler:     cliHandler,
		deliveryClient: deliveryClient,
	}, nil
}

func (a app) start() error {
	a.log.Info("приложение запускается")

	outPath := time.Now().Format("2006-01-02_15-04") + ".7z"

	err := a.cliHandler.Execute(os.Args[1:], outPath)
	if err != nil {
		return err
	}

	err = a.deliveryClient.Deliver(delivery.FileItem{
		Path: outPath,
		Name: outPath,
	})
	if err != nil {
		return err
	}

	ctx := context.Background()

	go a.runWorker(ctx, 1*time.Minute)

	select {}
}

// runWorker — бесконечный цикл плановой архивации и доставки.
func (a app) runWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			outPath := time.Now().Format("2006-01-02_15-04") + ".7z"

			if err := a.cliHandler.Execute([]string{}, outPath); err != nil {
				continue
			}

			err := a.deliveryClient.Deliver(delivery.FileItem{
				Path: outPath,
				Name: outPath,
			})
			if err != nil {
				a.log.Error("ошибка доставки", zap.Error(err))
			} else {
				a.log.Info("плановая задача успешно завершена")
			}
		}
	}
}
