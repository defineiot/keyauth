package assert

import (
	"strings"
	"testing"
)

func ErrorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func Contain(t *testing.T, err error, subString string) {
	if err == nil {
		t.Fatal("error is nil")
	} else if !strings.Contains(err.Error(), subString) {
		t.Fatalf("error %s have not contain %s", err.Error(), subString)
	}
}
