package file_test

import (
	"testing"

	"openauth/api/config/file"
)

func TestNewFileConf(t *testing.T) {
	fileconf := file.NewFileConf("../../../conf/openauth.conf")
	_, err := fileconf.GetConf()
	if err != nil {
		t.Fatal(err)
	}

}
