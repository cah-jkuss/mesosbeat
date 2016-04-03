package main

import (
	"github.com/kussj/mesosbeat/beater"

	"github.com/elastic/beats/libbeat/beat"
	"os"
)

var Name = "mesosbeat"

func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
