package liveupdates

import (
	"github.com/astro-walker/tilt/internal/store"
)

func HandleLiveUpdateUpsertAction(state *store.EngineState, action LiveUpdateUpsertAction) {
	n := action.LiveUpdate.Name
	state.LiveUpdates[n] = action.LiveUpdate
}

func HandleLiveUpdateDeleteAction(state *store.EngineState, action LiveUpdateDeleteAction) {
	delete(state.LiveUpdates, action.Name)
}
