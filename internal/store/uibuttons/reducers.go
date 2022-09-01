package uibuttons

import (
	"github.com/astro-walker/tilt/internal/store"
)

func HandleUIButtonUpsertAction(state *store.EngineState, action UIButtonUpsertAction) {
	n := action.UIButton.Name
	state.UIButtons[n] = action.UIButton
}

func HandleUIButtonDeleteAction(state *store.EngineState, action UIButtonDeleteAction) {
	delete(state.UIButtons, action.Name)
}
