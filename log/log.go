package log

import (
	"context"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	Log(ctx context.Context, level Level, keyvals ...interface{}) error
}

type LogHelper struct {
	logger Logger
}

func NewLogHelper(logger Logger) *LogHelper {
	return &LogHelper{
		logger: logger,
	}
}

func (l *LogHelper) Debug(ctx context.Context, keyvals ...interface{}) {
	l.logger.Log(ctx, DEBUG, keyvals...)
}

func (l *LogHelper) Info(ctx context.Context, keyvals ...interface{}) {
	l.logger.Log(ctx, INFO, keyvals...)
}

func (l *LogHelper) Warn(ctx context.Context, keyvals ...interface{}) {
	l.logger.Log(ctx, WARN, keyvals...)
}

func (l *LogHelper) Error(ctx context.Context, keyvals ...interface{}) {
	l.logger.Log(ctx, ERROR, keyvals...)
}

func (l *LogHelper) Fatal(ctx context.Context, keyvals ...interface{}) {
	l.logger.Log(ctx, FATAL, keyvals...)
	os.Exit(1)
}
