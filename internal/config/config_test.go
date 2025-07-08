package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	os.Chdir("../../dhall")
	defer os.Chdir(wd)

	if _, err := Load("./package.dhall"); err != nil {
		t.Error(err)
	}
}
