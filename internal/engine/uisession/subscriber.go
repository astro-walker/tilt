package uisession

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/astro-walker/tilt/internal/controllers/apicmp"
	"github.com/astro-walker/tilt/internal/hud/webview"
	"github.com/astro-walker/tilt/internal/store"
	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
)

type Subscriber struct {
	client ctrlclient.Client
}

func NewSubscriber(client ctrlclient.Client) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) OnChange(ctx context.Context, st store.RStore, summary store.ChangeSummary) error {
	if summary.IsLogOnly() {
		return nil
	}

	state := st.RLockState()
	session := webview.ToUISession(state)
	st.RUnlockState()

	stored := &v1alpha1.UISession{}
	err := s.client.Get(ctx, types.NamespacedName{Name: session.Name}, stored)
	if apierrors.IsNotFound(err) {
		// If nothing is stored, create it.
		err := s.client.Create(ctx, session)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		// If the cache hasn't started yet, that's OK.
		// We'll get it on the next OnChange()
		if _, ok := err.(*cache.ErrCacheNotStarted); ok {
			return nil
		}

		return err
	}

	if !apicmp.DeepEqual(session.Status, stored.Status) {
		// If the current version is different than what's stored, update it.
		update := &v1alpha1.UISession{
			ObjectMeta: *stored.ObjectMeta.DeepCopy(),
			Spec:       *stored.Spec.DeepCopy(),
			Status:     *session.Status.DeepCopy(),
		}
		err = s.client.Status().Update(ctx, update)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ store.Subscriber = &Subscriber{}
