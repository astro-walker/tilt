package prompt

import "github.com/astro-walker/tilt/internal/store"

type SwitchTerminalModeAction struct {
	Mode store.TerminalMode
}

func (SwitchTerminalModeAction) Action() {}

var _ store.Action = SwitchTerminalModeAction{}
