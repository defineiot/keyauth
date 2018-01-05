package tools_test

import (
	"testing"

	"openauth/tools"
)

func TestMakeUUID(t *testing.T) {
	uuid, err := tools.MakeUUID(24)
	if err != nil {
		t.Fatal(err)
	}

	if len(uuid) != 24 {
		t.Fatal("uuid lenghth not equal 24")
	}
}
