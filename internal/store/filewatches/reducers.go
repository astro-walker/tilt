package filewatches

import (
	"github.com/astro-walker/tilt/internal/store"
)

func HandleFileWatchUpsertAction(state *store.EngineState, action FileWatchUpsertAction) {
	n := action.FileWatch.Name
	state.FileWatches[n] = action.FileWatch
}

func HandleFileWatchDeleteAction(state *store.EngineState, action FileWatchDeleteAction) {
	delete(state.FileWatches, action.Name)
}
