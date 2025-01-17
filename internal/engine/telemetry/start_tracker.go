package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/astro-walker/tilt/internal/store"
)

type StartTracker struct {
	tracer        trace.Tracer
	span          trace.Span
	startFinished bool
}

func NewStartTracker(tracer trace.Tracer) *StartTracker {
	return &StartTracker{tracer: tracer, startFinished: false}
}

func (c *StartTracker) OnChange(ctx context.Context, st store.RStore, _ store.ChangeSummary) error {
	if c.startFinished {
		return nil
	}

	state := st.RLockState()
	defer st.RUnlockState()

	if !state.InitialBuildsCompleted() && c.span == nil {
		_, span := c.tracer.Start(ctx, "first_run")
		c.span = span
	}

	if state.InitialBuildsCompleted() && c.span != nil {
		c.span.End()
		c.startFinished = true
	}

	return nil
}
