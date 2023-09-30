package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	if _, err := Load("../../config.dhall"); err != nil {
		t.Error(err)
	}
}
