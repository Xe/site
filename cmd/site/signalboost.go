package main

import (
	"io/ioutil"

	"github.com/philandstuff/dhall-golang"
)

type Person struct {
	Name    string   `dhall:"name"`
	GitLink string   `dhall:"gitLink"`
	Twitter string   `dhall:"twitter"`
	Tags    []string `dhall:"tags"`
}

func loadPeople(path string) ([]Person, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var people []Person
	err = dhall.Unmarshal(data, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}
