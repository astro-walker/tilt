package cmd

import (
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/astro-walker/tilt/internal/controllers/apicmp"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

// Compares the exec-only fields of a CmdSpec.
// Ignores fields that specify dependency info (StartOn, RestartOn)
func cmdExecEqual(a, b v1alpha1.CmdSpec) bool {
	return execDelta.DeepEqual(a, b)
}

var execDelta = conversion.EqualitiesOrDie(
	append(
		[]interface{}{
			func(a, b *v1alpha1.StartOnSpec) bool { // ignore
				return true
			},
			func(a, b *v1alpha1.RestartOnSpec) bool { // ignore
				return true
			},
		},
		apicmp.Comparators()...)...)
