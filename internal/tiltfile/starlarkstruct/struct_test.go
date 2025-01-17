package starlarkstruct

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/astro-walker/tilt/internal/tiltfile/starkit"
)

func TestStruct(t *testing.T) {
	f := NewFixture(t)
	f.File("Tiltfile", `
x = struct(a = "foo", b = 2)
print("a",x.a)
print("b",x.b)
`)
	_, err := f.ExecFile("Tiltfile")
	require.NoError(t, err)
	require.Contains(t, f.PrintOutput(), "a foo")
	require.Contains(t, f.PrintOutput(), "b 2")
}

func TestModule(t *testing.T) {
	f := NewFixture(t)
	f.File("Tiltfile", `
x = module("test_module", a = "foo", b = 2)
print("a",x.a)
print("b",x.b)
print("x",x)
`)
	_, err := f.ExecFile("Tiltfile")
	require.NoError(t, err)
	require.Contains(t, f.PrintOutput(), "a foo")
	require.Contains(t, f.PrintOutput(), "b 2")
	require.Contains(t, f.PrintOutput(), "x <module \"test_module\">")
}

func NewFixture(tb testing.TB) *starkit.Fixture {
	return starkit.NewFixture(tb, NewPlugin())
}
