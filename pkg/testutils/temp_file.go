package testutils

import (
	"os"
	"io/ioutil"
)

type tempConfigFile struct {
	path string
}

func (t *tempConfigFile) Cleanup() {
	os.Remove(t.path)
}

func (t *tempConfigFile) GetPath() string {
	return t.path
}

func NewTempFile(content []byte) *tempConfigFile {
	f, _ := ioutil.TempFile("", "")
	f.Write(content)
	defer f.Close()
	return &tempConfigFile{
		path: f.Name(),
	}
}
