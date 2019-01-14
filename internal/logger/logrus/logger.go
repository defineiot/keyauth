package logrus

import (
	"fmt"
	"iot-pusher/app/log"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/logger/logrus/hooks"
)

// NewLogrusLogger use to get an logger instance
func NewLogrusLogger(opts *log.Opts) (logger.IOTAuthLogger, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate logger options error, %s", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logger := logrus.New()

	// set logger level
	level, err := logrus.ParseLevel(opts.Level)
	if err != nil {
		return nil, fmt.Errorf("parse logger level error, %s", err)
	}
	logger.SetLevel(level)

	// add context hook
	logger.AddHook(hooks.NewContextHook())
	// add file rotate hook
	if opts.FilePath != "" {
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

		logger.Hooks.Add(hookLF)
	}

	return logger, nil
}
