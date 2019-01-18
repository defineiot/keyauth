package logrus

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/logger/logrus/hooks"
)

// Logrus 基于logrus封装的日志实现
type Logrus struct {
	logger *logrus.Logger
}

// Trace 跟踪级别的日志
func (l *Logrus) Trace(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

// Debug debug 日志
func (l *Logrus) Debug(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info 通知日志
func (l *Logrus) Info(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn 警告
func (l *Logrus) Warn(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error 错误日志
func (l *Logrus) Error(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal 故障日志
func (l *Logrus) Fatal(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panic 崩溃日志
func (l *Logrus) Panic(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// WithFieldsTrace 额外信息
func (l *Logrus) WithFieldsTrace(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Printf(format, args...)
}

// WithFieldsDebug 额外信息
func (l *Logrus) WithFieldsDebug(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Debugf(format, args...)
}

// WithFieldsInfo 额外信息
func (l *Logrus) WithFieldsInfo(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Infof(format, args...)
}

// WithFieldsWarn 额外信息
func (l *Logrus) WithFieldsWarn(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Warnf(format, args...)
}

// WithFieldsError 额外信息
func (l *Logrus) WithFieldsError(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Errorf(format, args...)
}

// WithFieldsFatal 额外信息
func (l *Logrus) WithFieldsFatal(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Fatalf(format, args...)
}

// WithFieldsPanic 额外信息
func (l *Logrus) WithFieldsPanic(f logger.Fields, format string, args ...interface{}) {
	data := logrus.Fields(f)
	l.logger.WithFields(data).Panicf(format, args...)
}

// SetLevel 额外信息
func (l *Logrus) SetLevel(level logger.Level) {
	ll := logrus.Level(level)
	l.logger.SetLevel(ll)
}

// NewLogrusLogger use to get an logger instance
func NewLogrusLogger(opts *logger.Opts) (logger.Logger, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate logger options error, %s", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	log := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(opts.Level)
	if err != nil {
		return nil, fmt.Errorf("parse logger level error, %s", err)
	}
	log.SetLevel(level)

	if opts.NeedFileNumber {
		// 添加记录文件行号的日志插件
		log.AddHook(hooks.NewContextHook())
	}

	if opts.FilePath != "" {
		// 添加把日志记录到文件到插件
		tz, _ := time.LoadLocation("Local")
		writer, err := rotatelogs.New(
			opts.FilePath+".%Y-%m-%d-%H:%M",
			rotatelogs.WithLocation(tz),
			rotatelogs.WithLinkName(opts.FilePath),
			rotatelogs.WithMaxAge(time.Duration(604800)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
		)
		if err != nil {
			return nil, err
		}
		hookLF := lfshook.NewHook(lfshook.WriterMap{
			logrus.PanicLevel: writer,
			logrus.FatalLevel: writer,
			logrus.ErrorLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.InfoLevel:  writer,
			logrus.DebugLevel: writer,
		}, &logrus.JSONFormatter{})

		log.Hooks.Add(hookLF)
	}

	mylogger := &Logrus{logger: log}

	return mylogger, nil
}
