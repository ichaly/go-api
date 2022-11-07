package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var (
	logger    = newLogger()
	atomLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
)

// ResetLogger 初始化日志对象
func ResetLogger(level, filePath string) {
	atomLevel.SetLevel(transferLoglevel(level))
	logger = newLogger(filePath)
}

/*生成新的日志对象*/
func newLogger(filePath ...string) *zap.SugaredLogger {
	var (
		writer zapcore.WriteSyncer // 日志输出方
		// 原始的日志级别,可运行时修改,预置为debug
	)
	if len(filePath) != 0 && filePath[0] != "" {
		rotateHook := lumberjack.Logger{ // 日志滚动钩子
			Filename:   filePath[0], // 日志文件路径
			MaxSize:    1024,        // 文件最大多少M,默认100M
			MaxAge:     7,           // 最多保留多少天,默认不根据日期删除
			MaxBackups: 7,           // 最多保留多少个备份,
			LocalTime:  false,       // 是否使用本地时间,默认使用UTC时间
			Compress:   true,        // 是否gzip压缩
		}
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&rotateHook), os.Stdout)
	} else {
		writer = zapcore.NewMultiWriteSyncer(os.Stdout)
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(newEncoder()), writer, atomLevel)
	logger := zap.New(
		core,
		zap.AddCaller(),                       // 日志增加调用者
		zap.AddStacktrace(zapcore.PanicLevel), // panic级别下增加调用栈
		zap.AddCallerSkip(1),                  // 如果日志方法有封装,则调用者输出跳过的层数
	)
	return logger.Sugar()
}

// SetLevel 动态设置日志级别
func SetLevel(level string) {
	atomLevel.SetLevel(transferLoglevel(level))
}

// 定义了日志封装的各种设置
func newEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",                      // 输入信息的key名
		LevelKey:       "level",                        // 输出日志级别的key名
		TimeKey:        "time",                         // 输出时间的key名
		NameKey:        "logger",                       // 日志信息key名
		CallerKey:      "caller",                       // 调用者key名
		StacktraceKey:  "stacktrace",                   // 调用栈key名
		LineEnding:     zapcore.DefaultLineEnding,      // 每行分隔符,默认\n
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // level值的封装,配置为序列化为全大写
		EncodeTime:     timeFormatter,                  // 时间格式,配置为[2006-01-02 15:04:05]
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行消耗时间格式,配置为浮点秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 调用者格式,配置为包/文件:行号
		EncodeName:     zapcore.FullNameEncoder,        // 日志信息名处理,默认无处理
	}
}

/*将日志字符串级别转换为zap日志级别*/
func transferLoglevel(l string) zapcore.Level {
	l = strings.ToUpper(l)
	var level zapcore.Level
	switch l {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	case "PANIC":
		level = zapcore.PanicLevel
	case "FATAL":
		level = zapcore.FatalLevel
	default:
		level = zapcore.DebugLevel
	}
	return level
}

/*定义日志时间格式*/
func timeFormatter(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}
