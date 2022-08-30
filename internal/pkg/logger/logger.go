package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ILogger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message error, args ...any)
}

type Logger struct {
	logger *zap.Logger
}

func New() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.EncoderConfig.MessageKey = "message"

	logger, _ := config.Build()
	return &Logger{logger: logger}
}

func (l *Logger) Debug(message string, args ...any) {
	l.logger.Log(zap.DebugLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Info(message string, args ...any) {
	l.logger.Log(zap.InfoLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Warn(message string, args ...any) {
	l.logger.Log(zap.WarnLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Error(message error, args ...any) {
	l.logger.Log(zap.ErrorLevel, fmt.Sprintf(message.Error(), args...))
}

func (l *Logger) Exit() {
	l.logger.Sync()
}
