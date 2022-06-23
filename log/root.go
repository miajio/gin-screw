package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	param   *loggerParam // 日志参数
	paramMu sync.RWMutex // 日志参数锁

	logger *zap.SugaredLogger // 日志对象
	mu     sync.Mutex         // 同步锁
)

// Init 初始化日志参数
func Init(path string, maxSize, maxBackups, maxAge int, compress bool, logMap map[string]Level) {
	if param == nil {
		paramMu.Lock()
		defer paramMu.Unlock()
		if param == nil {
			param = &loggerParam{
				path:       path,
				maxSize:    maxSize,
				maxBackups: maxBackups,
				maxAge:     maxAge,
				compress:   compress,
				logMap:     logMap,
			}
		}
	}

	if logger == nil {
		mu.Lock()
		defer mu.Unlock()
		if logger == nil {
			logger = param.generateLogger()
		}
	}
}

// checkParam 检查日志对象
func checkLogger() {
	if logger == nil {
		panic("logger is nil; please call Init()")
	}
}

// GetLogger 获取日志对象
func GetLogger() *zap.SugaredLogger {
	checkLogger()
	return logger
}

// 生成日志对象
func (log *loggerParam) generateLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "log",
		CallerKey:     "lineNum",
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

	for name := range log.logMap {
		// logs[name]
		writer := GetWrite(fmt.Sprintf("%s/%s", log.path, name), log.maxSize, log.maxBackups, log.maxAge, log.compress)
		level := zap.LevelEnablerFunc(log.logMap[name])
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
