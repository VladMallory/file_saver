package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServiceName string
	Env         string
	LogLevel    string
}

func New(cfg Config) (*zap.Logger, error) {
	var zapConfig zap.Config
	var isJSON bool

	if cfg.Env == "json" || cfg.Env == "prod" {
		zapConfig = zap.NewProductionConfig()

		// Настройка формата времени
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		zapConfig.Encoding = "json"
		isJSON = true
	} else {
		zapConfig = zap.NewDevelopmentConfig()

		zapConfig.Encoding = "console"

		// Цветовая дифференциация уровней (INFO, WARN, ERROR)
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		// Привычный для человека формат времени (без лишних T и Z)
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

		// Сокращаем путь до вызвавшего файла (main.go:25 вместо /app/internal/...)
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		// Отключаем стектрейсы для INFO и WARN, оставляем только для ERROR и выше
		zapConfig.EncoderConfig.StacktraceKey = ""
	}

	// Для docker/k8s
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	// logger level
	if cfg.LogLevel != "" {
		level, err := zapcore.ParseLevel(cfg.LogLevel)
		if err == nil {
			zapConfig.Level = zap.NewAtomicLevelAt(level)
		}
	}

	// Создание логгера
	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	if isJSON {
		logger = logger.With(
			zap.String("service", cfg.ServiceName),
			zap.String("env", cfg.Env),
		)
	}

	return logger, nil
}
