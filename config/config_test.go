package config

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	if _, err := Parse("./config.ts"); err != nil {
		t.Fatal(err)
	}
}

func TestDate(t *testing.T) {
	t.Run("unmarshal", func(t *testing.T) {
		d := Date{}
		if err := d.UnmarshalJSON([]byte("\"2020-01-02\"")); err != nil {
			t.Fatal(err)
		}
		if d.Year() != 2020 || d.Month() != 1 || d.Day() != 2 {
			t.Fatal("wrong date")
		}
	})

	t.Run("marshal", func(t *testing.T) {
		d := Date{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)}
		b, err := d.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "\"2020-01-02\"" {
			t.Fatal("wrong date")
		}
	})
}
