package server

import (
	"github.com/astro-walker/tilt/pkg/model"
)

type AppendToTriggerQueueAction struct {
	Name   model.ManifestName
	Reason model.BuildReason
}

func (AppendToTriggerQueueAction) Action() {}

// TODO: a way to clear an override
type OverrideTriggerModeAction struct {
	ManifestNames []model.ManifestName
	TriggerMode   model.TriggerMode
}

func (OverrideTriggerModeAction) Action() {}
