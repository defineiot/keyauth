package testutils

import "os"

type envManager struct {
	recodes []string
}

func (e *envManager) Cleanup() {
	for _, recode := range e.recodes {
		os.Unsetenv(recode)
	}
}

func (e *envManager) Set(key, value string) {
	if err := os.Setenv(key, value); err == nil {
		e.recodes = append(e.recodes, key)
	}
}

func (e *envManager) Unset(key string) {
	var recodes []string
	if err := os.Unsetenv(key); err == nil {
		for _, r := range e.recodes {
			if r != key {
				recodes = append(recodes, r)
			}
		}
		e.recodes = recodes
	}
}

func NewEnvManager() *envManager {
	return &envManager{}
}
