package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogLevel  = zapcore.DebugLevel
	defaultLogFile   = "log.json"
	defaultFileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)

// ILogger абстракция для логирования.
type ILogger interface {
	// Debug сообщение уровня отладки.
	Debug(message string, args ...any)
	// Info информационное сообщение.
	Info(message string, args ...any)
	// Warn предупреждение.
	Warn(message string, args ...any)
	// Error сообщение об ошибке.
	Error(message error, args ...any)
	// Fatal выводит сообщение и вызывает os.Exit(1).
	Fatal(message error, args ...any)
}

// Logger реализация логгера (ILogger).
type Logger struct {
	logger *zap.Logger
}

// New создаёт объект Logger
func New() *Logger {
	config := newEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	//fileEncoder := zapcore.NewJSONEncoder(config)
	//logFile, _ := os.OpenFile(defaultLogFile, defaultFileFlags, 0644)
	//writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		//zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{logger: logger}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 | 15:04:05.9999"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
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

func (l *Logger) Fatal(message error, args ...any) {
	l.logger.Log(zap.FatalLevel, fmt.Sprintf(message.Error(), args...))
}

func (l *Logger) Exit() {
	l.logger.Sync()
}
