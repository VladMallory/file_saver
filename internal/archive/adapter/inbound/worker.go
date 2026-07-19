package archiveinboud

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type BackupUseCase interface {
	Run() error
}

type Scheduler struct {
	log      *zap.Logger
	uc       BackupUseCase
	interval time.Duration
}

func NewScheduler(log *zap.Logger, uc BackupUseCase, interval time.Duration) *Scheduler {
	return &Scheduler{
		log:      log,
		uc:       uc,
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.log.Info(
		"фоновый планировщик запущен",
		zap.Duration("interval", s.interval),
	)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("планировщик остановлен по сигналу")
			return

		case <-ticker.C:
			start := time.Now()

			if err := s.uc.Run(); err != nil {
				s.log.Error(
					"ошибка выполения плановой задачи",
					zap.Error(err),
				)
			} else {
				s.log.Info(
					"плановая задача успешно завершена",
					zap.Duration("duration",
						time.Since(start)),
				)
			}
		}
	}
}
