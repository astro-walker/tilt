package dockerimages

import (
	"github.com/astro-walker/tilt/internal/store"
)

func HandleDockerImageUpsertAction(state *store.EngineState, action DockerImageUpsertAction) {
	obj := action.DockerImage
	n := obj.Name
	state.DockerImages[n] = obj
}

func HandleDockerImageDeleteAction(state *store.EngineState, action DockerImageDeleteAction) {
	delete(state.DockerImages, action.Name)
}
