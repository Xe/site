package main

import (
	"log"

	"github.com/google/gops/agent"
)

func init() {
	if err := agent.Listen(nil); err != nil {
		log.Fatal(err)
	}
}
