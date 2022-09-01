package tiltfile

import (
	"github.com/astro-walker/tilt/internal/container"
	"github.com/astro-walker/tilt/internal/k8s/testyaml"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
	"github.com/astro-walker/tilt/pkg/model"
)

const SanchoDockerfile = `
FROM go:1.10
ADD . .
RUN go install github.com/tilt-dev/sancho
ENTRYPOINT /go/bin/sancho
`

var SanchoRef = container.MustParseSelector(testyaml.SanchoImage)

type pathFixture interface {
	Path() string
}

func NewSanchoDockerBuildImageTarget(f pathFixture) model.ImageTarget {
	return model.MustNewImageTarget(SanchoRef).
		WithDockerImage(v1alpha1.DockerImageSpec{
			DockerfileContents: SanchoDockerfile,
			Context:            f.Path(),
		})
}
