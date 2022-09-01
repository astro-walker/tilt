package buildcontrols

import (
	"context"

	"github.com/astro-walker/tilt/internal/ospath"
	"github.com/astro-walker/tilt/pkg/logger"
	"github.com/astro-walker/tilt/pkg/model"
)

type BuildEntry struct {
	Name         model.ManifestName
	BuildReason  model.BuildReason
	FilesChanged []string
}

func LogBuildEntry(ctx context.Context, entry BuildEntry) {
	buildReason := entry.BuildReason
	changedFiles := entry.FilesChanged
	firstBuild := buildReason.Has(model.BuildReasonFlagInit)

	l := logger.Get(ctx).WithFields(logger.Fields{logger.FieldNameBuildEvent: "init"})
	if firstBuild {
		l.Infof("Initial Build")
	} else {
		if len(changedFiles) > 0 {
			t := "File"
			if len(changedFiles) > 1 {
				t = "Files"
			}
			l.Infof("%d %s Changed: %s", len(changedFiles), t, ospath.FormatFileChangeList(changedFiles))
		} else {
			l.Infof("%s", buildReason)
		}
	}
}
