package tiltfile

import (
	"github.com/google/wire"

	"github.com/astro-walker/tilt/internal/tiltfile/config"
	"github.com/astro-walker/tilt/internal/tiltfile/k8scontext"
	"github.com/astro-walker/tilt/internal/tiltfile/tiltextension"
	"github.com/astro-walker/tilt/internal/tiltfile/version"
)

var WireSet = wire.NewSet(
	ProvideTiltfileLoader,
	k8scontext.NewPlugin,
	version.NewPlugin,
	config.NewPlugin,
	tiltextension.NewPlugin,
)
