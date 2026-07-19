package deliverytelegram

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	delivery "saveFile/internal/deliveryArchive/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type TelegramSender struct {
	log    *zap.Logger
	bot    *bot.Bot
	chatID string
}

func NewTelegramSender(log *zap.Logger, token, chatID string) (*TelegramSender, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}

	return &TelegramSender{
		log:    log,
		bot:    b,
		chatID: chatID,
	}, nil
}

// Send открывает файл по пути из FileItem и отправляет его в Telegram.
func (t *TelegramSender) Send(file delivery.FileItem) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer closeHelper(f, &err)

	params := &bot.SendDocumentParams{
		ChatID: t.chatID,
		Document: &models.InputFileUpload{
			Filename: filepath.Base(file.Path),
			Data:     f,
		},
	}

	if _, err := t.bot.SendDocument(context.Background(), params); err != nil {
		t.log.Error("ошибка отправки файла в Telegram", zap.Error(err))
		return err
	}

	return nil
}

func closeHelper(closer io.Closer, err *error) {
	if cerr := closer.Close(); cerr != nil {
		*err = errors.Join(*err, cerr)
	}
}
