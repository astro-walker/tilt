package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

var (
	sel = MustParseSelector("gcr.io/foo/bar")
)

func TestNewRefSetWithInvalidRegistryErrors(t *testing.T) {
	reg := &v1alpha1.RegistryHosting{Host: "invalid"}
	assertNewRefSetError(t, sel, reg, "repository name must be canonical")
}

func TestNewRefSetErrorsWithBadLocalRef(t *testing.T) {
	// Force "repository name must not be longer than 255 characters" when assembling LocalRef
	var longname string
	for i := 0; i < 230; i++ {
		longname += "o"
	}
	selector := MustParseSelector(longname)
	reg := &v1alpha1.RegistryHosting{Host: "gcr.io/somewhat/long/hostname"}
	assertNewRefSetError(t, selector, reg, "after applying default registry")
}

func TestNewRefSetErrorsWithBadClusterRef(t *testing.T) {
	// Force "repository name must not be longer than 255 characters" when assembling ClusterRef
	var longname string
	for i := 0; i < 230; i++ {
		longname += "o"
	}
	selector := MustParseSelector(longname)
	reg := &v1alpha1.RegistryHosting{Host: "gcr.io", HostFromContainerRuntime: "gcr.io/somewhat/long/hostname"}
	assertNewRefSetError(t, selector, reg, "after applying default registry")
}

func TestNewRefSetEmptyRegistryOK(t *testing.T) {
	_, err := NewRefSet(sel, nil)
	assert.NoError(t, err)
}

var cases = []struct {
	name               string
	host               string
	clusterHost        string
	expectedLocalRef   string
	expectedClusterRef string
}{
	{"empty registry", "", "", "gcr.io/foo", "gcr.io/foo"},
	{"host only", "localhost:1234", "", "localhost:1234/gcr.io_foo", "localhost:1234/gcr.io_foo"},
	{"host and clusterHost", "localhost:1234", "registry:1234", "localhost:1234/gcr.io_foo", "registry:1234/gcr.io_foo"},
}

func TestDeriveRefs(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reg *v1alpha1.RegistryHosting
			if tc.host != "" {
				reg = &v1alpha1.RegistryHosting{
					Host:                     tc.host,
					HostFromContainerRuntime: tc.clusterHost,
				}
			}
			refs, err := NewRefSet(MustParseSelector("gcr.io/foo"), reg)
			require.NoError(t, err)

			localRef := refs.LocalRef()
			clusterRef := refs.ClusterRef()

			assert.Equal(t, tc.expectedLocalRef, localRef.String())
			assert.Equal(t, tc.expectedClusterRef, clusterRef.String())
		})
	}
}

func TestWithoutRegistry(t *testing.T) {
	reg := &v1alpha1.RegistryHosting{
		Host:                     "localhost:5000",
		HostFromContainerRuntime: "localhost:5000",
	}
	refs, err := NewRefSet(MustParseSelector("foo"), reg)
	require.NoError(t, err)

	assert.Equal(t, "localhost:5000/foo", FamiliarString(refs.LocalRef()))
	assert.Equal(t, "foo", FamiliarString(refs.WithoutRegistry().LocalRef()))
}

func assertNewRefSetError(t *testing.T, selector RefSelector, reg *v1alpha1.RegistryHosting, expectedErr string) {
	t.Helper()
	_, err := NewRefSet(selector, reg)
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr)
}
