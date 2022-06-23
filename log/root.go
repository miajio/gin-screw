package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerParam 日志参数
type LoggerParam struct {
	path       string // 日志文件路径
	maxSize    int    // 日志最大存储量
	maxBackups int    // 日志最大备份数
	maxAge     int    // 日志最大存储天数
	compress   bool   // 是否压缩
}

/*
NewLogger 创建日志参数对象
@param path			日志文件路径
@param maxSize		日志最大存储量
@param maxBackups	日志最大备份数
@param maxAge		日志最大存储天数
@param compress		是否压缩
@return 日志参数对象
*/
func NewLogger(path string, maxSize, maxBackups, maxAge int, compress bool) *LoggerParam {
	return &LoggerParam{
		path:       path,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

/*
New 创建Uber日志对象
@param logs 日志文件
@return uber日志对象
*/
func (log *LoggerParam) New(logs map[string]Level) *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "log",
		CallerKey:     "lineNo",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(time.Format("[2006-01-02 15:04:05]"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	// 	return level < zapcore.WarnLevel && level >= zap.InfoLevel
	// })

	cores := make([]zapcore.Core, 0)

	for name := range logs {
		// logs[name]
		writer := GetWrite(fmt.Sprintf("%s/%s", log.path, name), log.maxSize, log.maxBackups, log.maxAge, log.compress)
		level := zap.LevelEnablerFunc(logs[name])
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(writer), level))
	}

	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel))

	core := zapcore.NewTee(
		// zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(log.getWrite(log.InfoName)), infoLevel),
		// zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(log.getWrite(log.ErrorName)), errorLevel),
		// json log
		// zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel),
		// console log
		cores...,
	)

	caller := zap.AddCaller()

	development := zap.Development()
	return zap.New(core, caller, development, zap.Fields()).Sugar()
}

func GetWrite(path string, maxSize, maxBackups, maxAge int, compress bool) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,       // 文件路径
		MaxSize:    maxSize,    // 日志最大存储量
		MaxBackups: maxBackups, // 日志最大备份数
		MaxAge:     maxAge,     // 日志最大存储天数
		Compress:   compress,   // 是否压缩
	}
}
