package logrus

import (
	"fmt"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"openauth/api/logger"
	"openauth/api/logger/logrus/hooks"
)

// NewLogrusLogger use to get an logger instance
func NewLogrusLogger(opts *logger.Opts) (logger.OpenAuthLogger, error) {
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
		})

		logger.Hooks.Add(hookLF)
	}

	return logger, nil
}
