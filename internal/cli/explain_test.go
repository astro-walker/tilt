package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

func TestExplain(t *testing.T) {
	f := newServerFixture(t)

	err := f.client.Create(f.ctx, &v1alpha1.Cmd{
		ObjectMeta: metav1.ObjectMeta{Name: "my-sleep"},
		Spec: v1alpha1.CmdSpec{
			Args: []string{"sleep", "1"},
		},
	})
	require.NoError(t, err)

	out := bytes.NewBuffer(nil)
	streams := genericclioptions.IOStreams{Out: out}

	explain := newExplainCmd(streams)
	explain.register()

	err = explain.run(f.ctx, []string{"cmd"})
	require.NoError(t, err)

	assert.Contains(t, out.String(), `Cmd represents a process on the host machine.`)
}
