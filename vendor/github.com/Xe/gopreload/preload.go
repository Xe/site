//+build linux,go1.8

package gopreload

import (
	"log"
	"os"
	"plugin"
	"strings"
)

func init() {
	gpv := os.Getenv("GO_PRELOAD")
	if gpv == "" {
		return
	}

	for _, elem := range strings.Split(gpv, ",") {
		log.Printf("gopreload: trying to open: %s", elem)
		_, err := plugin.Open(elem)
		if err != nil {
			log.Printf("%v from GO_PRELOAD cannot be loaded: %v", elem, err)
			continue
		}
	}
}
