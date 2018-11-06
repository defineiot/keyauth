package file_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/conf/file"
)

func TestNewFileConf(t *testing.T) {
	fileconf := file.NewFileConf("../../../.keyauth/keyauth.conf")
	_, err := fileconf.GetConf()
	if err != nil {
		t.Fatal(err)
	}

}
