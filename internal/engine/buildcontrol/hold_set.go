package buildcontrol

import (
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/pkg/model"
)

type HoldSet map[model.ManifestName]store.Hold

func (s HoldSet) IsEligible(target *store.ManifestTarget) bool {
	mn := target.Manifest.Name
	if s[mn].Reason != store.HoldReasonNone {
		return false
	}

	if target.State.IsBuilding() {
		return false
	}

	if target.NextBuildReason() == 0 {
		return false
	}

	return true
}

func (s HoldSet) RemoveIneligibleTargets(targets []*store.ManifestTarget) []*store.ManifestTarget {
	result := make([]*store.ManifestTarget, 0, len(targets))
	for _, target := range targets {
		if s.IsEligible(target) {
			result = append(result, target)
		}
	}
	return result
}

func (s HoldSet) AddHold(target *store.ManifestTarget, hold store.Hold) {
	mn := target.Manifest.Name
	if s[mn].Reason != store.HoldReasonNone {
		return
	}

	if target.State.IsBuilding() {
		return
	}

	if target.NextBuildReason() == 0 {
		return
	}

	s[mn] = hold
}

// For all the targets that should have built and don't have a prior Hold, add the given Hold.
func (s HoldSet) Fill(targets []*store.ManifestTarget, hold store.Hold) {
	for _, target := range targets {
		s.AddHold(target, hold)
	}
}
