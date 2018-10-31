package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type locale map[string]map[string]string

func (l locale) Value(group, key string, args ...interface{}) string {
	sg, ok := l[group]
	if !ok {
		return "no group " + group
	}

	result, ok := sg[key]
	if !ok {
		return fmt.Sprintf("in group %s, no key %s", group, key)
	}

	return fmt.Sprintf(result, args...)
}

type translations struct {
	locales map[string]locale
}

func (t *translations) LoadLocale(name string, r io.Reader) error {
	l := locale{}
	d := json.NewDecoder(r)
	err := d.Decode(&l)
	if err == nil {
		t.locales[name] = l
	}
	return err
}

func (t *translations) Get(name string) (locale, bool) {
	l, ok := t.locales[name]
	if !ok {
		return nil, false
	}

	return l, ok
}
