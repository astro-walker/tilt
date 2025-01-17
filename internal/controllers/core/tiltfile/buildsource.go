package tiltfile

import (
	"context"

	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/pkg/model"
	"github.com/astro-walker/tilt/pkg/model/logstore"
)

// BuildEntry is vestigial, but currently used to help manage state about a tiltfile build.
type BuildEntry struct {
	Name                  model.ManifestName
	FilesChanged          []string
	BuildReason           model.BuildReason
	Args                  []string
	TiltfilePath          string
	CheckpointAtExecStart logstore.Checkpoint
	LoadCount             int
	ArgsChanged           bool
}

func (be *BuildEntry) WithLogger(ctx context.Context, st store.RStore) context.Context {
	return store.WithManifestLogHandler(ctx, st, be.Name, SpanIDForLoadCount(be.Name, be.LoadCount))
}
