package main

import (
	"context"
	"log"
	"os/signal"
	"saveFile/internal/config"
	"saveFile/pkg/logger"
	"syscall"
	"time"

	deliverytelegram "saveFile/internal/deliveryArchive/adapter/outbound"
	delivery "saveFile/internal/deliveryArchive/domain"
	deliveryusecase "saveFile/internal/deliveryArchive/service"

	"go.uber.org/zap"
)

type app struct {
	logClient      *zap.Logger
	deliveryClient deliveryusecase.DeliveryService
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.run()
	if err != nil {
		log.Fatal(err)
	}
}

func newApp() (app, error) {
	cfg, err := config.Load()
	if err != nil {
		return app{}, err
	}

	// ===logger===
	logClient, err := logger.New(logger.Config{
		ServiceName: "save-file-service",
		Env:         cfg.Env,
		LogLevel:    cfg.LogLevel,
	})
	if err != nil {
		return app{}, err
	}

	// ===архивация===
	// pathArchiveClient := patharchive.NewPathProvider(log)
	// sevenzipClient := sevenzip.NewArchiver(log)
	// archiveClient := archiveusecase.NewArchiveService(log, pathArchiveClient, sevenzipClient)

	// ===доставка===
	telegramClient, err := deliverytelegram.NewTelegramSender(
		logClient, cfg.TelegramToken,
		cfg.TelegramChatID,
	)
	if err != nil {
		return app{}, err
	}

	deliveryClient := deliveryusecase.NewDeliveryService(logClient, telegramClient)

	logClient.Info("приложение собралось")

	return app{
		logClient:      logClient,
		deliveryClient: deliveryClient,
	}, nil
}

func (a app) run() error {
	a.logClient.Info("приложение запускается")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	outPath := time.Now().Format("2006-01-02_15-04") + ".7z"

	err := a.deliveryClient.Deliver(delivery.FileItem{
		Path: outPath,
		Name: outPath,
	})
	if err != nil {
		return err
	}

	go a.runWorker(ctx, 1*time.Minute)

	<-ctx.Done()
	a.logClient.Info("приложение завершается")

	return nil
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

			err := a.deliveryClient.Deliver(delivery.FileItem{
				Path: outPath,
				Name: outPath,
			})
			if err != nil {
				a.logClient.Error("ошибка доставки", zap.Error(err))
			} else {
				a.logClient.Info("плановая задача успешно завершена")
			}
		}
	}
}
