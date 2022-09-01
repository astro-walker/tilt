package liveupdate

import (
	"github.com/astro-walker/tilt/internal/build"
	"github.com/astro-walker/tilt/internal/store/liveupdates"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Input struct {
	// Derived from DockerResource
	IsDC bool

	// Derived from KubernetesResource + KubenetesSelector + DockerResource
	Containers []liveupdates.Container

	// Derived from FileWatch + Sync rules
	ChangedFiles []build.PathMapping

	LastFileTimeSynced metav1.MicroTime
}
