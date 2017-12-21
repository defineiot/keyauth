package assert

import (
	"fmt"
	"strings"
	"testing"
	"runtime"
	"path/filepath"
)

func ErrorIsNil(t *testing.T, err error) {
	if err != nil {
		who := resourceCall()
		t.Fatalf("%s %s", who, err.Error())
	}
}

func Contain(t *testing.T, err error, subString string) {
	who := resourceCall()
	if err == nil {
		t.Fatalf("%s error is nil", who)
	} else if !strings.Contains(err.Error(), subString) {
		t.Fatalf("%s error %s have not contain %s", who, err.Error(), subString)
	}
}

func Equal(t *testing.T, actual, expect string) {
	if actual != expect {
		who := resourceCall()
		t.Fatalf("%s %s not equal %s", who, actual, expect)
	}
}

func resourceCall() string {
	_, filename, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d: ", filepath.Base(filename), line)
}
