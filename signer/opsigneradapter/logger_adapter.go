package opsigneradapter

import (
	"context"
	"log/slog"

	signercommon "github.com/agglayer/go_signer/common"
	"github.com/ethereum/go-ethereum/log"
)

type LoggerAdapter struct {
	logger signercommon.Logger
}

// Check that LoggerAdapter implements log.Logger
var _ log.Logger = (*LoggerAdapter)(nil)

func NewLoggerAdapter(logger signercommon.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: logger,
	}
}

func (l *LoggerAdapter) Crit(msg string, ctx ...interface{}) {
	l.Log(slog.LevelError, msg, ctx...)
}

func (l *LoggerAdapter) Debug(msg string, ctx ...interface{}) {
	l.Log(slog.LevelDebug, msg, ctx...)
}

func (l *LoggerAdapter) Trace(msg string, ctx ...interface{}) {
	l.Log(slog.LevelDebug, msg, ctx...)
}
func (l *LoggerAdapter) Info(msg string, ctx ...interface{}) {
	l.Log(slog.LevelInfo, msg, ctx...)
}

func (l *LoggerAdapter) Error(msg string, ctx ...interface{}) {
	l.Log(slog.LevelError, msg, ctx...)
}

func (l *LoggerAdapter) Warn(msg string, ctx ...interface{}) {
	l.Log(slog.LevelWarn, msg, ctx...)
}

func (l *LoggerAdapter) Log(level slog.Level, msg string, ctx ...interface{}) {
	if level >= slog.LevelError {
		l.logger.Error(append([]interface{}{msg}, ctx...)...)
	} else if level >= slog.LevelWarn {
		l.logger.Warn(append([]interface{}{msg}, ctx...)...)
	} else if level >= slog.LevelInfo {
		l.logger.Info(append([]interface{}{msg}, ctx...)...)
	}
	l.logger.Debug(append([]interface{}{msg}, ctx...)...)
}

func (l *LoggerAdapter) Write(level slog.Level, msg string, attrs ...any) {
	l.Log(level, msg, attrs...)
}

func (l *LoggerAdapter) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (l *LoggerAdapter) Handler() slog.Handler {
	return nil
}

func (l *LoggerAdapter) With(ctx ...interface{}) log.Logger {
	return nil
}

func (l *LoggerAdapter) New(ctx ...interface{}) log.Logger {
	return nil
}
