package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astro-walker/tilt/internal/tiltfile/starkit"
)

func TestMetricsEnabled(t *testing.T) {
	f := newFixture(t)
	f.File("Tiltfile", "experimental_metrics_settings(enabled=True)")
	_, err := f.ExecFile("Tiltfile")
	assert.NoError(t, err)
	assert.Contains(t, f.PrintOutput(), "deprecated")
}

func newFixture(tb testing.TB) *starkit.Fixture {
	return starkit.NewFixture(tb, NewPlugin())
}
