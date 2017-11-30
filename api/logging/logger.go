package logging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
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

func (o *Opts) validate() error {
	if o.Name == "" {
		return errors.New("the logger name missed")
	}
	if o.Level == "" {
		o.Level = "info"
	}

	if o.FilePath != "" {
		var err error
		// get the absolute path
		if !filepath.IsAbs(o.FilePath) {
			o.FilePath, err = filepath.Abs(o.FilePath)
			if err != nil {
				return fmt.Errorf("get the file absolute path error, %s", err)
			}
		}
		// if file not exits make all
		dirp := filepath.Dir(o.FilePath)
		_, err = os.Stat(dirp)
		if err != nil && os.IsNotExist(err) {
			if err := os.MkdirAll(dirp, os.ModePerm); err != nil {
				return fmt.Errorf("mkdir for log path error, %s", err)
			}
		}
	}

	return nil
}
