package main

import (
	"context"
	"time"

	delivery "saveFile/internal/deliveryArchive/domain"

	"go.uber.org/zap"
)

// runWorker - простой метод структуры app, который крутит бесконечный цикл
func (a app) runWorker(ctx context.Context, interval time.Duration) {
	a.log.Info("воркер запущен", zap.Duration("interval", interval))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.log.Info("воркер остановлен")
			return

		case <-ticker.C:
			a.log.Info("запуск плановой задачи")

			outPath := time.Now().Format("2006-01-02_15-04") + ".7z"

			if err := a.cliHandler.Execute([]string{}, outPath); err != nil {
				a.log.Error("ошибка архивации", zap.Error(err))
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
