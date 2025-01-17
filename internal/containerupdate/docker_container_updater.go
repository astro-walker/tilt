package containerupdate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/astro-walker/tilt/internal/container"
	"github.com/astro-walker/tilt/internal/docker"
	"github.com/astro-walker/tilt/internal/k8s"
	"github.com/astro-walker/tilt/internal/store/liveupdates"
	"github.com/astro-walker/tilt/pkg/logger"
	"github.com/astro-walker/tilt/pkg/model"
)

type DockerUpdater struct {
	dCli docker.Client
}

var _ ContainerUpdater = &DockerUpdater{}

func NewDockerUpdater(dCli docker.Client) *DockerUpdater {
	return &DockerUpdater{dCli: dCli}
}

func (cu *DockerUpdater) WillBuildToKubeContext(kctx k8s.KubeContext) bool {
	return cu.dCli.Env().WillBuildToKubeContext(kctx)
}

func (cu *DockerUpdater) UpdateContainer(ctx context.Context, cInfo liveupdates.Container,
	archiveToCopy io.Reader, filesToDelete []string, cmds []model.Cmd, hotReload bool) error {
	l := logger.Get(ctx)

	err := cu.rmPathsFromContainer(ctx, cInfo.ContainerID, filesToDelete)
	if err != nil {
		return errors.Wrap(err, "rmPathsFromContainer")
	}

	// Use `tar` to unpack the files into the container.
	//
	// Although docker has a copy API, it's buggy and not well-maintained
	// (whereas the Exec API is part of the CRI and much more battle-tested).
	// Discussion:
	// https://github.com/astro-walker/tilt/issues/3708
	tarCmd := tarCmd()
	err = cu.dCli.ExecInContainer(ctx, cInfo.ContainerID, tarCmd, archiveToCopy, l.Writer(logger.InfoLvl))
	if err != nil {
		if exitCode, ok := ExtractExitCode(err); ok {
			return wrapTarExecErr(err, tarCmd, exitCode)
		}
		return fmt.Errorf("copying changed files: %w", err)
	}

	// Exec run's on container
	for i, cmd := range cmds {
		l.Infof("[CMD %d/%d] %s", i+1, len(cmds), strings.Join(cmd.Argv, " "))
		err = cu.dCli.ExecInContainer(ctx, cInfo.ContainerID, cmd, nil, l.Writer(logger.InfoLvl))
		if err != nil {
			return fmt.Errorf(
				"executing on container %s: %w",
				cInfo.ContainerID.ShortStr(),
				wrapRunStepError(wrapDockerGenericExecErr(cmd, err)),
			)
		}
	}

	if hotReload {
		l.Debugf("Hot reload on, skipping container restart: %s", cInfo.DisplayName())
		return nil
	}

	// Restart container so that entrypoint restarts with the updated files etc.
	l.Debugf("Restarting container: %s", cInfo.DisplayName())
	err = cu.dCli.ContainerRestartNoWait(ctx, cInfo.ContainerID.String())
	if err != nil {
		return errors.Wrap(err, "ContainerRestart")
	}
	return nil
}

func (cu *DockerUpdater) rmPathsFromContainer(ctx context.Context, cID container.ID, paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	out := bytes.NewBuffer(nil)
	err := cu.dCli.ExecInContainer(ctx, cID, model.Cmd{Argv: makeRmCmd(paths)}, nil, out)
	if err != nil {
		if docker.IsExitError(err) {
			return fmt.Errorf("Error deleting files from container: %s", out.String())
		}
		return errors.Wrap(err, "Error deleting files from container")
	}
	return nil
}

func makeRmCmd(paths []string) []string {
	cmd := []string{"rm", "-rf"}
	cmd = append(cmd, paths...)
	return cmd
}

func wrapDockerGenericExecErr(cmd model.Cmd, err error) error {
	if exitCode, ok := ExtractExitCode(err); ok {
		return NewExecError(cmd, exitCode)
	}
	return err
}
