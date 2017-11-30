package logging

import (
	"errors"
	"fmt"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"openauth/api/logging/hooks"
)

// Opts is an options for get logger
type Opts struct {
	Name     string
	Level    string
	FilePath string
}

// NewLogger use to get an logger instance
func NewLogger(opts *Opts) (*logrus.Logger, error) {
	if err := opts.validate(); err != nil {
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

	// add hooks
	logger.AddHook(hooks.NewHook())
	if opts.FilePath != "" {
		logger.AddHook(lfshook.NewHook(lfshook.PathMap{
			logrus.DebugLevel: opts.FilePath,
			logrus.InfoLevel:  opts.FilePath,
			logrus.WarnLevel:  opts.FilePath,
			logrus.ErrorLevel: opts.FilePath,
			logrus.FatalLevel: opts.FilePath,
			logrus.PanicLevel: opts.FilePath,
		}))
	}

	return logger, nil
}

func (o *Opts) validate() error {
	if o.Name == "" {
		return errors.New("the logger name missed")
	}
	if o.Level == "" {
		o.Level = "info"
	}

	return nil
}
