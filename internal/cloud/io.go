package cloud

import (
	"context"
	"io"
	"os"

	"google.golang.org/protobuf/types/known/timestamppb"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/astro-walker/tilt/internal/hud/webview"
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/pkg/logger"
	proto_webview "github.com/astro-walker/tilt/pkg/webview"
)

type Snapshotter struct {
	st     store.RStore
	client ctrlclient.Client
}

func NewSnapshotter(st store.RStore, client ctrlclient.Client) *Snapshotter {
	return &Snapshotter{
		st:     st,
		client: client,
	}
}

func (s *Snapshotter) WriteSnapshot(ctx context.Context, path string) {
	view, err := webview.CompleteView(ctx, s.client, s.st)
	if err != nil {
		logger.Get(ctx).Errorf("Fetching snapshot: %v", err)
		return
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Get(ctx).Errorf("Writing snapshot to file: %v", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	snapshot := &proto_webview.Snapshot{
		View:      view,
		CreatedAt: timestamppb.Now(),
	}

	err = WriteSnapshotTo(ctx, snapshot, f)
	if err != nil {
		logger.Get(ctx).Errorf("Writing snapshot to file: %v", err)
		return
	}
}

func WriteSnapshotTo(ctx context.Context, snapshot *proto_webview.Snapshot, w io.Writer) error {
	jsEncoder := &runtime.JSONPb{
		OrigName: false,
		Indent:   "  ",
	}
	return jsEncoder.NewEncoder(w).Encode(snapshot)
}
