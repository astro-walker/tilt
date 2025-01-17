package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

func TestApply(t *testing.T) {
	f := newServerFixture(t)

	f.WriteFile("sleep.yaml", `
apiVersion: tilt.dev/v1alpha1
kind: Cmd
metadata:
  name: my-sleep
spec:
  args: ["sleep", "1"]
`)
	out := bytes.NewBuffer(nil)
	streams := genericclioptions.IOStreams{Out: out}

	cmd := newApplyCmd(streams)
	c := cmd.register()
	err := c.Flags().Parse([]string{"-f", f.JoinPath("sleep.yaml")})
	require.NoError(t, err)

	err = cmd.run(f.ctx, c.Flags().Args())
	require.NoError(t, err)
	assert.Contains(t, out.String(), `cmd.tilt.dev/my-sleep created`)

	var sleep v1alpha1.Cmd
	err = f.client.Get(f.ctx, types.NamespacedName{Name: "my-sleep"}, &sleep)
	require.NoError(t, err)
	assert.Equal(t, []string{"sleep", "1"}, sleep.Spec.Args)
}
