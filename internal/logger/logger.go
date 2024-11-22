package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Приватный глобальный объект логгера
var globalLogger *zap.Logger

// Init создает zap логгер
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// Паблик обертки для доступа к методам глобального объекта логгера

// Debug обертка над zap Debug
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info обертка над zap Info
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn обертка над zap Warn
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error обертка над zap Error
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal обертка над zap Fatal
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// WithOptions обертка над zap WithOptions
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
