package main

import "testing"

func TestLoadPeople(t *testing.T) {
	people, err := loadPeople("../../signalboost.dhall")
	if err != nil {t.Fatal(err)}

	for _, person := range people {
		t.Run(person.Name, func(t *testing.T) {
			if person.Name == "" {
				t.Error("missing name")
			}

			if len(person.Tags) == 0 {
				t.Error("missing tags")
			}

			if person.Twitter == "" {
				t.Error("missing twitter")
			}

			if person.GitLink == "" {
				t.Error("missing git link")
			}
		})
	}
}
