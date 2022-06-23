package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level func(level zapcore.Level) bool

var (
	DebugLevel = func(level zapcore.Level) bool {
		return level == zap.DebugLevel
	}

	InfoLevel = func(level zapcore.Level) bool {
		return level == zap.InfoLevel
	}

	ErrorLevel = func(level zapcore.Level) bool {
		return level > zap.InfoLevel
	}
)
