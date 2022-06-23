package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level func(level zapcore.Level) bool

// loggerParam 日志参数
type loggerParam struct {
	path       string           // 日志文件路径
	maxSize    int              // 日志最大存储量
	maxBackups int              // 日志最大备份数
	maxAge     int              // 日志最大存储天数
	compress   bool             // 是否压缩
	logMap     map[string]Level // 日志级别字典
}

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
