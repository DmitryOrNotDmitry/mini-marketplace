package logger

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLog *zap.SugaredLogger

	once = new(sync.Once)
)

// LoggerConfig содержит параметры для логгера
type LoggerConfig struct {
	Level       zapcore.Level
	ServiceName string
}

// InitLogger инициализирует глобальный логгер
func InitLogger(logConfig *LoggerConfig) {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
		config.Level.SetLevel(logConfig.Level)

		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		l, err := config.Build(
			zap.AddCallerSkip(1),
			zap.Fields(zap.String("service", logConfig.ServiceName)),
		)
		if err != nil {
			panic("failed to build logger")
		}
		globalLog = l.Sugar()
	})
}

func getLogger() *zap.SugaredLogger {
	return globalLog
}

// Warnw выводит сообщение об предупреждении в лог
func Warnw(msg string, keysAndValues ...interface{}) {
	getLogger().Warnw(msg, keysAndValues...)
}

// WarnwCtx выводит сообщение об предупреждении в лог, добавляя данные из контекста
func WarnwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = addTracingVars(ctx, keysAndValues)
	getLogger().Warnw(msg, keysAndValues...)
}

// Errorw выводит сообщение об ошибке в лог
func Errorw(msg string, keysAndValues ...interface{}) {
	getLogger().Errorw(msg, keysAndValues...)
}

// ErrorwCtx выводит сообщение об ошибке в лог, добавляя данные из контекста
func ErrorwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = addTracingVars(ctx, keysAndValues)
	getLogger().Errorw(msg, keysAndValues...)
}

// Infow выводит информационное сообщение в лог
func Infow(msg string, keysAndValues ...interface{}) {
	getLogger().Infow(msg, keysAndValues...)
}

// InfowCtx выводит информационное сообщение в лог, добавляя данные из контекста
func InfowCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = addTracingVars(ctx, keysAndValues)
	getLogger().Infow(msg, keysAndValues...)
}

func addTracingVars(ctx context.Context, keysAndValues []interface{}) []interface{} {
	span := trace.SpanFromContext(ctx)
	if sc := span.SpanContext(); sc.IsValid() {
		keysAndValues = append(keysAndValues,
			"trace_id", sc.TraceID().String(),
			"span_id", sc.SpanID().String(),
		)
	}
	return keysAndValues
}

// Sync отправляет логи из буффера.
func Sync() error {
	return getLogger().Sync()
}
