package cluster

import (
	"github.com/google/wire"

	"github.com/astro-walker/tilt/internal/controllers/apis/cluster"
)

var WireSet = wire.NewSet(
	NewConnectionManager,
	wire.Bind(new(cluster.ClientProvider), new(*ConnectionManager)),
	wire.InterfaceValue(new(KubernetesClientFactory), KubernetesClientFunc(KubernetesClientFromEnv)),
	wire.InterfaceValue(new(DockerClientFactory), DockerClientFunc(DockerClientFromEnv)),
)
