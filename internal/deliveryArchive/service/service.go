package deliveryusecase

import (
	"os"

	delivery "saveFile/internal/deliveryArchive/domain"

	"go.uber.org/zap"
)

type Sender interface {
	Send(file delivery.FileItem) error
}

type DeliveryService struct {
	log    *zap.Logger
	sender Sender
}

func NewDeliveryService(log *zap.Logger, sender Sender) DeliveryService {
	return DeliveryService{
		log:    log,
		sender: sender,
	}
}

func (s DeliveryService) Deliver(file delivery.FileItem) error {
	if _, err := os.Stat(file.Path); err != nil {
		return delivery.ErrNoFile
	}

	if err := s.sender.Send(file); err != nil {
		s.log.Error(
			"доставка не удалась",
			zap.String("file", file.Path),
			zap.Error(err),
		)

		return delivery.ErrSendFailed
	}

	s.log.Info(
		"файл успешно доставлен",
		zap.String("file", file.Path),
		zap.String("name", file.Name),
	)

	return nil
}
