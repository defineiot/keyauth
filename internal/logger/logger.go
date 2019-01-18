package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// 定义各种日志级别, 总共7个级别
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

// Level 定义日志级别
type Level int

// Fields 日志元数据
type Fields map[string]interface{}

// Opts is an options for get logger
type Opts struct {
	Name           string
	Level          string
	FilePath       string
	NeedFileNumber bool
}

// Logger is openauth's logger
type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})

	WithFieldsTrace(f Fields, format string, args ...interface{})
	WithFieldsDebug(f Fields, format string, args ...interface{})
	WithFieldsInfo(f Fields, format string, args ...interface{})
	WithFieldsWarn(f Fields, format string, args ...interface{})
	WithFieldsError(f Fields, format string, args ...interface{})
	WithFieldsFatal(f Fields, format string, args ...interface{})
	WithFieldsPanic(f Fields, format string, args ...interface{})

	SetLevel(level Level)
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
		// 获取配置的日志文件存放目录的绝对路径
		if !filepath.IsAbs(o.FilePath) {
			o.FilePath, err = filepath.Abs(o.FilePath)
			if err != nil {
				return fmt.Errorf("get the file absolute path error, %s", err)
			}
		}

		// 如果配置的日志文件不存在就创建目录
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
