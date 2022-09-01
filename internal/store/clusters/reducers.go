package clusters

import (
	"github.com/astro-walker/tilt/internal/store"
)

func HandleClusterUpsertAction(state *store.EngineState, action ClusterUpsertAction) {
	n := action.Cluster.Name
	state.Clusters[n] = action.Cluster
}

func HandleClusterDeleteAction(state *store.EngineState, action ClusterDeleteAction) {
	delete(state.Clusters, action.Name)
}
