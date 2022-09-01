package telemetry

import (
	"github.com/astro-walker/tilt/pkg/model"
)

type TelemetryScriptRanAction struct {
	Status model.TelemetryStatus
}

func (TelemetryScriptRanAction) Action() {}
