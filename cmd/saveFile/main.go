package main

import (
	"context"
	"log"
	"os"
	"time"

	archiveinboud "saveFile/internal/archive/adapter/inbound"
	"saveFile/internal/archive/adapter/outbound/archiveformat/sevenzip"
	patharchive "saveFile/internal/archive/adapter/outbound/path"
	archiveusecase "saveFile/internal/archive/service"
	"saveFile/internal/config"
	deliverytelegram "saveFile/internal/deliveryArchive/adapter"
	delivery "saveFile/internal/deliveryArchive/domain"
	deliveryusecase "saveFile/internal/deliveryArchive/service"
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
		if app.log != nil {
			_ = app.log.Sync()
		}
	}()

	err = app.run()
	if err != nil {
		log.Fatal(err)
	}
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

	//===доставка===
	telegramClient, err := deliverytelegram.NewTelegramSender(log, cfg.TelegramToken, cfg.TelegramChatID)
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

func (a app) run() error {
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
