package main

import (
	_ "expvar"

	"github.com/astro-walker/tilt/internal/cli"
	"github.com/astro-walker/tilt/pkg/model"
)

// Magic variables set by goreleaser
var version string
var date string

func main() {
	cli.SetTiltInfo(model.TiltBuild{
		Version: version,
		Date:    date,
	})
	cli.Execute()
}
