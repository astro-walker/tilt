package manifestutils

import (
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
	"github.com/astro-walker/tilt/pkg/model"
)

func NewManifestTargetWithPod(m model.Manifest, pod v1alpha1.Pod) *store.ManifestTarget {
	mt := store.NewManifestTarget(m)
	mt.State.RuntimeState = store.NewK8sRuntimeStateWithPods(m, pod)
	return mt
}
