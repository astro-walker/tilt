package kubernetesapplys

import (
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/internal/store/kubernetesdiscoverys"
)

func HandleKubernetesApplyUpsertAction(state *store.EngineState, action KubernetesApplyUpsertAction) {
	n := action.KubernetesApply.Name
	state.KubernetesApplys[n] = action.KubernetesApply
	kubernetesdiscoverys.RefreshKubernetesResource(state, n)
}

func HandleKubernetesApplyDeleteAction(state *store.EngineState, action KubernetesApplyDeleteAction) {
	delete(state.KubernetesApplys, action.Name)
	kubernetesdiscoverys.RefreshKubernetesResource(state, action.Name)
}
