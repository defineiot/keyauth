package hooks_test

import (
	"bytes"
	"testing"

	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"

	"openauth/api/logger/logrus/hooks"
)

func TestNewContextHook(t *testing.T) {
	contextH := hooks.NewContextHook()
	log := logrus.New()

	log.AddHook(contextH)
	log.SetLevel(logrus.DebugLevel)
	log.Formatter = &logrus.JSONFormatter{}

	buffer := new(bytes.Buffer)
	log.Out = buffer

	log.Debug("error")

	source := jsoniter.Get(buffer.Bytes(), "source").ToString()
	if source != "hooks/context_test.go:24" {
		t.Fatal("context debug plugin context find error")
	}

}
