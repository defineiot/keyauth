package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// IOTAuthLogger is openauth's logger
type IOTAuthLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

// Opts is an options for get logger
type Opts struct {
	Name     string
	Level    string
	FilePath string
}

// Validate the logger config
func (o *Opts) Validate() error {
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
