//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package buildcontrol

import (
	"context"

	"github.com/google/wire"
	"github.com/jonboulle/clockwork"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tilt-dev/wmclient/pkg/dirs"

	"github.com/tilt-dev/clusterid"
	"github.com/astro-walker/tilt/internal/analytics"
	"github.com/astro-walker/tilt/internal/build"
	"github.com/astro-walker/tilt/internal/containerupdate"
	"github.com/astro-walker/tilt/internal/controllers/core/cmdimage"
	"github.com/astro-walker/tilt/internal/controllers/core/dockercomposeservice"
	"github.com/astro-walker/tilt/internal/controllers/core/dockerimage"
	"github.com/astro-walker/tilt/internal/controllers/core/kubernetesapply"
	"github.com/astro-walker/tilt/internal/docker"
	"github.com/astro-walker/tilt/internal/dockercompose"
	"github.com/astro-walker/tilt/internal/dockerfile"
	"github.com/astro-walker/tilt/internal/k8s"
	"github.com/astro-walker/tilt/internal/localexec"
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/internal/store/liveupdates"
	"github.com/astro-walker/tilt/internal/tracer"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

var BaseWireSet = wire.NewSet(
	// dockerImageBuilder ( = ImageBuilder)
	wire.Value(dockerfile.Labels{}),

	v1alpha1.NewScheme,
	k8s.ProvideMinikubeClient,
	build.NewDockerBuilder,
	build.NewCustomBuilder,
	wire.Bind(new(build.DockerKubeConnection), new(*build.DockerBuilder)),

	// BuildOrder
	NewDockerComposeBuildAndDeployer,
	NewImageBuildAndDeployer,
	NewLocalTargetBuildAndDeployer,
	containerupdate.NewDockerUpdater,
	containerupdate.NewExecUpdater,
	build.NewImageBuilder,

	tracer.InitOpenTelemetry,

	liveupdates.ProvideUpdateMode,
)

func ProvideImageBuildAndDeployer(
	ctx context.Context,
	docker docker.Client,
	kClient k8s.Client,
	env clusterid.Product,
	kubeContext k8s.KubeContext,
	clusterEnv docker.ClusterEnv,
	dir *dirs.TiltDevDir,
	clock build.Clock,
	kp build.KINDLoader,
	analytics *analytics.TiltAnalytics,
	ctrlclient ctrlclient.Client,
	st store.RStore,
	execer localexec.Execer) (*ImageBuildAndDeployer, error) {
	wire.Build(
		BaseWireSet,
		kubernetesapply.NewReconciler,
		dockerimage.NewReconciler,
		cmdimage.NewReconciler,
	)

	return nil, nil
}

func ProvideDockerComposeBuildAndDeployer(
	ctx context.Context,
	dcCli dockercompose.DockerComposeClient,
	dCli docker.Client,
	ctrlclient ctrlclient.Client,
	st store.RStore,
	clock clockwork.Clock,
	dir *dirs.TiltDevDir) (*DockerComposeBuildAndDeployer, error) {
	wire.Build(
		BaseWireSet,
		dockercomposeservice.WireSet,
		build.ProvideClock,
		build.NewKINDLoader,
		dockerimage.NewReconciler,
		cmdimage.NewReconciler,
	)

	return nil, nil
}
