package dockerprune

import (
	"fmt"
	"time"

	"go.starlark.net/starlark"

	"github.com/astro-walker/tilt/pkg/model"

	"github.com/astro-walker/tilt/internal/tiltfile/starkit"
)

// Implements functions for dealing with Docker Prune settings.
type Plugin struct {
}

func NewPlugin() Plugin {
	return Plugin{}
}

func (e Plugin) NewState() interface{} {
	return model.DockerPruneSettings{
		Enabled:    true,
		MaxAge:     model.DockerPruneDefaultMaxAge,
		Interval:   model.DockerPruneDefaultInterval,
		KeepRecent: model.DockerPruneDefaultKeepRecent,
	}
}

func (e Plugin) OnStart(env *starkit.Environment) error {
	return env.AddBuiltin("docker_prune_settings", e.dockerPruneSettings)
}

func (e Plugin) dockerPruneSettings(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var disable bool
	var keepRecent starlark.Value
	var intervalHrs, numBuilds, maxAgeMins int
	if err := starkit.UnpackArgs(thread, fn.Name(), args, kwargs,
		"disable?", &disable,
		"max_age_mins?", &maxAgeMins,
		"num_builds?", &numBuilds,
		"interval_hrs?", &intervalHrs,
		"keep_recent?", &keepRecent); err != nil {
		return nil, err
	}

	if numBuilds != 0 && intervalHrs != 0 {
		return nil, fmt.Errorf("can't specify both 'prune every X builds' and 'prune every Y hours'; please pass " +
			"only one of `num_builds` and `interval_hrs`")
	}

	err := starkit.SetState(thread, func(settings model.DockerPruneSettings) (model.DockerPruneSettings, error) {
		settings.Enabled = !disable
		if maxAgeMins != 0 {
			settings.MaxAge = time.Duration(maxAgeMins) * time.Minute
		}
		settings.NumBuilds = numBuilds
		if intervalHrs != 0 {
			settings.Interval = time.Duration(intervalHrs) * time.Hour
		}
		if keepRecent != nil {
			recent, err := starlark.AsInt32(keepRecent)
			if err != nil {
				return settings, err
			}
			settings.KeepRecent = recent
		}
		return settings, nil
	})

	return starlark.None, err
}

var _ starkit.StatefulPlugin = Plugin{}

func MustState(model starkit.Model) model.DockerPruneSettings {
	state, err := GetState(model)
	if err != nil {
		panic(err)
	}
	return state
}

func GetState(m starkit.Model) (model.DockerPruneSettings, error) {
	var state model.DockerPruneSettings
	err := m.Load(&state)
	return state, err
}
