package tools_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/tools"
)

func TestMakeUUID(t *testing.T) {
	uuid := tools.MakeUUID(24)

	if len(uuid) != 24 {
		t.Fatal("uuid lenghth not equal 24")
	}
}
