package tiltfiles

import (
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

type TiltfileUpsertAction struct {
	Tiltfile *v1alpha1.Tiltfile
}

func NewTiltfileUpsertAction(tf *v1alpha1.Tiltfile) TiltfileUpsertAction {
	return TiltfileUpsertAction{Tiltfile: tf.DeepCopy()}
}

func (TiltfileUpsertAction) Action() {}

type TiltfileDeleteAction struct {
	Name string
}

func NewTiltfileDeleteAction(n string) TiltfileDeleteAction {
	return TiltfileDeleteAction{Name: n}
}

func (TiltfileDeleteAction) Action() {}
