package containerupdate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/util/exec"

	"github.com/astro-walker/tilt/internal/build"
	"github.com/astro-walker/tilt/internal/sliceutils"

	"github.com/astro-walker/tilt/internal/k8s"
	"github.com/astro-walker/tilt/internal/testutils"
	"github.com/astro-walker/tilt/pkg/model"
)

var toDelete = []string{"/foo/delete_me", "/bar/me_too"}
var (
	cmdA = model.Cmd{Argv: []string{"a"}}
	cmdB = model.Cmd{Argv: []string{"b", "bar", "baz"}}
)
var cmds = []model.Cmd{cmdA, cmdB}

func TestUpdateContainerDoesntSupportRestart(t *testing.T) {
	f := newExecFixture(t)

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("boop"), toDelete, cmds, false)
	if assert.NotNil(t, err, "expect Exec UpdateContainer to fail if !hotReload") {
		assert.Contains(t, err.Error(), "ExecUpdater does not support `restart_container()` step")
	}
}

func TestUpdateContainerDeletesFiles(t *testing.T) {
	f := newExecFixture(t)

	// No files to delete
	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("boop"), nil, cmds, true)
	if err != nil {
		t.Fatal(err)
	}

	for _, call := range f.kCli.ExecCalls {
		if sliceutils.StringSliceStartsWith(call.Cmd, "rm") {
			t.Fatal("found kubernetes exec `rm` call, expected none b/c no files to delete")
		}
	}

	// Two files to delete
	err = f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("boop"), toDelete, cmds, true)
	if err != nil {
		t.Fatal(err)
	}
	var rmCmd []string
	for _, call := range f.kCli.ExecCalls {
		if sliceutils.StringSliceStartsWith(call.Cmd, "rm") {
			if len(rmCmd) != 0 {
				t.Fatalf(`found two rm commands, expected one.
cmd 1: %v
cmd 2: %v`, rmCmd, call.Cmd)
			}
			rmCmd = call.Cmd
		}
	}
	if len(rmCmd) == 0 {
		t.Fatal("no `rm` cmd found, expected one b/c we specified files to delete")
	}

	expectedRmCmd := []string{"rm", "-rf", "/foo/delete_me", "/bar/me_too"}
	assert.Equal(t, expectedRmCmd, rmCmd)
}

func TestUpdateContainerTarsArchive(t *testing.T) {
	f := newExecFixture(t)

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("hello world"), nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	expectedCmd := []string{"tar", "-C", "/", "-x", "-f", "-"}
	if assert.Len(t, f.kCli.ExecCalls, 1, "expect exactly 1 k8s exec call") {
		call := f.kCli.ExecCalls[0]
		assert.Equal(t, expectedCmd, call.Cmd)
		assert.Equal(t, []byte("hello world"), call.Stdin)
	}
}

func TestUpdateContainerRunsCommands(t *testing.T) {
	f := newExecFixture(t)

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("hello world"), nil, cmds, true)
	if err != nil {
		t.Fatal(err)
	}

	if assert.Len(t, f.kCli.ExecCalls, 3, "expect exactly 3 k8s exec calls") {
		// second and third calls should be our cmd runs
		assert.Equal(t, cmdA.Argv, f.kCli.ExecCalls[1].Cmd)
		assert.Equal(t, cmdB.Argv, f.kCli.ExecCalls[2].Cmd)
	}
}

func TestUpdateContainerRunsFailure(t *testing.T) {
	f := newExecFixture(t)

	// The first exec() call is a copy, so won't trigger a RunStepFailure
	f.kCli.ExecErrors = []error{
		nil,
		exec.CodeExitError{Err: fmt.Errorf("Compile error"), Code: 1234},
	}

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("hello world"), nil, cmds, true)
	if assert.True(t, build.IsRunStepFailure(err)) {
		assert.Equal(t, `executing on container test_conta: command "a" failed with exit code: 1234`, err.Error())
	}
	assert.Equal(t, 2, len(f.kCli.ExecCalls))
}

func TestUpdateContainerMissingTarFailure(t *testing.T) {
	f := newExecFixture(t)

	f.kCli.ExecErrors = []error{
		errors.New("opaque Kubernetes error that includes the phrase 'executable file not found' in it"),
	}

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("hello world"), nil, cmds, true)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Please check that the container image includes `tar` in $PATH.")
	}
	assert.Equal(t, 1, len(f.kCli.ExecCalls))
}

func TestUpdateContainerPermissionDenied(t *testing.T) {
	f := newExecFixture(t)

	f.kCli.ExecOutputs = []io.Reader{strings.NewReader("tar: app/index.js: Cannot open: File exists\n")}
	f.kCli.ExecErrors = []error{exec.CodeExitError{Err: fmt.Errorf("command terminated with exit code 2"), Code: 2}}

	err := f.ecu.UpdateContainer(f.ctx, TestContainerInfo, newReader("hello world"), nil, cmds, true)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "container filesystem denied access")
	}
	assert.Equal(t, 1, len(f.kCli.ExecCalls))
}

type execUpdaterFixture struct {
	t    testing.TB
	ctx  context.Context
	kCli *k8s.FakeK8sClient
	ecu  *ExecUpdater
}

func newExecFixture(t testing.TB) *execUpdaterFixture {
	fakeCli := k8s.NewFakeK8sClient(t)
	cu := &ExecUpdater{
		kCli: fakeCli,
	}
	ctx, _, _ := testutils.CtxAndAnalyticsForTest()

	return &execUpdaterFixture{
		t:    t,
		ctx:  ctx,
		kCli: fakeCli,
		ecu:  cu,
	}
}

func newReader(contents string) io.Reader {
	return bytes.NewBuffer([]byte(contents))
}
