package cloud

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/astro-walker/tilt/internal/hud/webview"
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/internal/store/tiltfiles"
	"github.com/astro-walker/tilt/internal/testutils"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
	"github.com/astro-walker/tilt/pkg/model"
	proto_webview "github.com/astro-walker/tilt/pkg/webview"
)

func TestWriteSnapshotTo(t *testing.T) {
	ctx, _, _ := testutils.CtxAndAnalyticsForTest()
	buf := bytes.NewBuffer(nil)

	state := store.NewState()
	tiltfiles.HandleTiltfileUpsertAction(state, tiltfiles.TiltfileUpsertAction{
		Tiltfile: &v1alpha1.Tiltfile{
			ObjectMeta: metav1.ObjectMeta{Name: model.MainTiltfileManifestName.String()},
			Spec:       v1alpha1.TiltfileSpec{Path: "Tiltfile"},
		},
	})
	now := time.Unix(1551202573, 0)
	snapshot := &proto_webview.Snapshot{
		View: &proto_webview.View{
			UiSession: webview.ToUISession(*state),
		},
		CreatedAt: timestamppb.New(now),
	}

	resources, err := webview.ToUIResourceList(*state, make(map[string][]v1alpha1.DisableSource))
	require.NoError(t, err)
	snapshot.View.UiResources = resources

	for _, r := range resources {
		for i, cond := range r.Status.Conditions {
			// Clear the transition timestamps so that the test is hermetic.
			cond.LastTransitionTime = metav1.MicroTime{}
			r.Status.Conditions[i] = cond
		}
	}

	startTime := timestamppb.New(now)
	snapshot.View.TiltStartTime = startTime

	err = WriteSnapshotTo(ctx, snapshot, buf)
	assert.NoError(t, err)
	assert.Equal(t, `{
  "view": {
    "tiltStartTime": "2019-02-26T17:36:13Z",
    "uiSession": {
      "metadata": {
        "name": "Tiltfile"
      },
      "status": {
        "versionSettings": {
          "checkUpdates": true
        },
        "tiltfileKey": "Tiltfile"
      }
    },
    "uiResources": [
      {
        "metadata": {
          "name": "(Tiltfile)"
        },
        "status": {
          "runtimeStatus": "not_applicable",
          "updateStatus": "pending",
          "order": 1,
          "conditions": [
            {
              "type": "UpToDate",
              "status": "False",
              "reason": "UpdatePending"
            },
            {
              "type": "Ready",
              "status": "False",
              "reason": "UpdatePending"
            }
          ]
        }
      }
    ]
  },
  "createdAt": "2019-02-26T17:36:13Z"
}
`, buf.String())
}
