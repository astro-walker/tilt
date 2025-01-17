package watch

import (
	"go.starlark.net/starlark"

	"github.com/astro-walker/tilt/internal/tiltfile/starkit"
	"github.com/astro-walker/tilt/internal/tiltfile/value"
	"github.com/astro-walker/tilt/pkg/model"
)

type Plugin struct {
}

func NewPlugin() Plugin {
	return Plugin{}
}

func (e Plugin) NewState() interface{} {
	return model.WatchSettings{}
}

func (e Plugin) OnStart(env *starkit.Environment) error {
	return env.AddBuiltin("watch_settings", e.setWatchSettings)
}

func (e Plugin) setWatchSettings(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	err := starkit.SetState(thread, func(settings model.WatchSettings) (model.WatchSettings, error) {
		var ignores value.StringOrStringList
		if err := starkit.UnpackArgs(thread, fn.Name(), args, kwargs,
			"ignore?", &ignores,
		); err != nil {
			return settings, err
		}

		if len(ignores.Values) != 0 {
			settings.Ignores = append(settings.Ignores, model.Dockerignore{
				LocalPath: starkit.AbsWorkingDir(thread),
				Patterns:  ignores.Values,
				Source:    "watch_settings()",
			})
		}

		return settings, nil
	})

	return starlark.None, err
}

var _ starkit.StatefulPlugin = Plugin{}

func MustState(model starkit.Model) model.WatchSettings {
	state, err := GetState(model)
	if err != nil {
		panic(err)
	}
	return state
}

func GetState(m starkit.Model) (model.WatchSettings, error) {
	var state model.WatchSettings
	err := m.Load(&state)
	return state, err
}
