package logger

import (
	"onbio/conf"

	log "onbio/zaplog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger is a run log instance
	Logger *zap.Logger
)

func Init() {
	logConfig := conf.GetLogConfig()

	cfg := log.Config{
		EncodeLogsAsJson:   true,
		FileLoggingEnabled: true,
		Directory:          logConfig.FilePath,
		Filename:           logConfig.FileName,
		MaxSize:            512,
		MaxBackups:         30,
		MaxAge:             7,
	}
	Logger = log.NewLogger(cfg)
}

// Debug func
func Debug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, fields...)
}

// Info func
func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}

// Warn func
func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, fields...)
}

// Error func
func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}

// Panic func
func Panic(msg string, fields ...zapcore.Field) {
	Logger.Panic(msg, fields...)
}

// Fatal func
func Fatal(msg string, fields ...zapcore.Field) {
	Logger.Fatal(msg, fields...)
}
