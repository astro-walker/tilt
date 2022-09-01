package kubernetesdiscovery

import "github.com/astro-walker/tilt/internal/store"

type Dispatcher interface {
	Dispatch(action store.Action)
}
