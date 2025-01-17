package kubernetesapply

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/astro-walker/tilt/internal/k8s/testyaml"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

// The hashes are hard-coded in this file to ensure we
// don't accidentally change them.
//
// When updating the hashes, make sure that you don't accidentally
// change two hashes to the same value
func MustComputeInputHash(t testing.TB, spec v1alpha1.KubernetesApplySpec, imageMaps map[types.NamespacedName]*v1alpha1.ImageMap) string {
	hash, err := ComputeInputHash(spec, imageMaps)
	require.NoError(t, err)
	return hash
}

func TestComputeHashSancho(t *testing.T) {
	spec := v1alpha1.KubernetesApplySpec{YAML: testyaml.SanchoYAML}
	hash := MustComputeInputHash(t, spec, nil)
	assert.Equal(t, hash, "elgAsMuqc3wosKnB-vrr3mmMDPc=")
}

func TestComputeHashSanchoSidecar(t *testing.T) {
	spec := v1alpha1.KubernetesApplySpec{YAML: testyaml.SanchoSidecarYAML}
	hash := MustComputeInputHash(t, spec, nil)
	assert.Equal(t, hash, "SE9Jjlp-yakBXqPs2IjnrhCous8=")
}

func TestComputeHashSanchoImageMap(t *testing.T) {
	spec := v1alpha1.KubernetesApplySpec{YAML: testyaml.SanchoYAML, ImageMaps: []string{"sancho"}}
	imageMaps := make(map[types.NamespacedName]*v1alpha1.ImageMap)
	imageMaps[types.NamespacedName{Name: "sancho"}] = &v1alpha1.ImageMap{
		ObjectMeta: metav1.ObjectMeta{Name: "sancho"},
		Spec:       v1alpha1.ImageMapSpec{Selector: "sancho"},
		Status:     v1alpha1.ImageMapStatus{Image: "sancho:1234"},
	}

	hash := MustComputeInputHash(t, spec, imageMaps)
	assert.Equal(t, hash, "xwdU6LXhS59jR4lfTApIkP26RJ0=")
}

func TestComputeHashSanchoIgnoresIrrelevantImageMap(t *testing.T) {
	spec := v1alpha1.KubernetesApplySpec{YAML: testyaml.SanchoYAML}
	imageMaps := make(map[types.NamespacedName]*v1alpha1.ImageMap)
	imageMaps[types.NamespacedName{Name: "sancho"}] = &v1alpha1.ImageMap{
		ObjectMeta: metav1.ObjectMeta{Name: "sancho"},
		Spec:       v1alpha1.ImageMapSpec{Selector: "sancho"},
		Status:     v1alpha1.ImageMapStatus{Image: "sancho:1234"},
	}

	assert.Equal(t, MustComputeInputHash(t, spec, nil), MustComputeInputHash(t, spec, imageMaps))
}

func TestComputeHashSanchoImageMapChange(t *testing.T) {
	spec := v1alpha1.KubernetesApplySpec{YAML: testyaml.SanchoYAML, ImageMaps: []string{"sancho"}}
	imageMaps := make(map[types.NamespacedName]*v1alpha1.ImageMap)
	imageMaps[types.NamespacedName{Name: "sancho"}] = &v1alpha1.ImageMap{
		ObjectMeta: metav1.ObjectMeta{Name: "sancho"},
		Spec:       v1alpha1.ImageMapSpec{Selector: "sancho"},
		Status:     v1alpha1.ImageMapStatus{Image: "sancho:45678"},
	}

	hash := MustComputeInputHash(t, spec, imageMaps)
	assert.Equal(t, hash, "VzBtR4akztTBxoGNlWc56QtD9b0=")
}
