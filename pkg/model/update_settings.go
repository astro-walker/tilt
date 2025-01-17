package model

import (
	"time"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

const (
	DefaultMaxParallelUpdates = 3
)

type UpdateSettings struct {
	maxParallelUpdates int           // max number of updates to run concurrently
	k8sUpsertTimeout   time.Duration // timeout for k8s upsert operations

	// A list of images to suppress the warning for.
	SuppressUnusedImageWarnings []string
}

func (us UpdateSettings) MaxParallelUpdates() int {
	// Min. value is 1
	if us.maxParallelUpdates < 1 {
		return 1
	}
	return us.maxParallelUpdates
}

func (us UpdateSettings) WithMaxParallelUpdates(n int) UpdateSettings {
	// Min. value is 1
	if n < 1 {
		n = 1
	}
	us.maxParallelUpdates = n
	return us
}

func (us UpdateSettings) K8sUpsertTimeout() time.Duration {
	// Min. value is 1s
	if us.k8sUpsertTimeout < time.Second {
		return time.Second
	}
	return us.k8sUpsertTimeout
}

func (us UpdateSettings) WithK8sUpsertTimeout(timeout time.Duration) UpdateSettings {
	// Min. value is 1s
	if us.k8sUpsertTimeout < time.Second {
		timeout = time.Second
	}
	us.k8sUpsertTimeout = timeout
	return us
}

func DefaultUpdateSettings() UpdateSettings {
	return UpdateSettings{
		maxParallelUpdates: DefaultMaxParallelUpdates,
		k8sUpsertTimeout:   v1alpha1.KubernetesApplyTimeoutDefault,
	}
}
