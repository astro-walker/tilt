package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

func TestCmdExecEqual(t *testing.T) {
	assert.True(t,
		cmdExecEqual(
			v1alpha1.CmdSpec{Args: []string{"cat"}},
			v1alpha1.CmdSpec{Args: []string{"cat"}}))
	assert.False(t,
		cmdExecEqual(
			v1alpha1.CmdSpec{Args: []string{"cat"}},
			v1alpha1.CmdSpec{Args: []string{"dog"}}))
	assert.True(t,
		cmdExecEqual(
			v1alpha1.CmdSpec{Args: []string{"cat"}},
			v1alpha1.CmdSpec{
				Args:      []string{"cat"},
				StartOn:   &v1alpha1.StartOnSpec{UIButtons: []string{"x"}},
				RestartOn: &v1alpha1.RestartOnSpec{FileWatches: []string{"x"}},
			}))
}
