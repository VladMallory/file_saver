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

			if err := a.cliHandler.Execute([]string{}); err != nil {
				a.log.Error("ошибка архивации", zap.Error(err))
				continue
			}

			// 2. Вызываем доставку
			err := a.deliveryClient.Deliver(delivery.FileItem{
				Path: "backup.7z",
				Name: "backup.7z",
			})
			if err != nil {
				a.log.Error("ошибка доставки", zap.Error(err))
			} else {
				a.log.Info("плановая задача успешно завершена")
			}
		}
	}
}
