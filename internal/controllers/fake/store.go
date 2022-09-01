package fake

import (
	"io"

	"github.com/astro-walker/tilt/internal/store"
)

type testStore struct {
	*store.TestingStore
	out io.Writer
}

func NewTestingStore(out io.Writer) *testStore {
	return &testStore{
		TestingStore: store.NewTestingStore(),
		out:          out,
	}
}

func (s *testStore) Dispatch(action store.Action) {
	s.TestingStore.Dispatch(action)
	if action, ok := action.(store.LogAction); ok {
		_, _ = s.out.Write(action.Message())
	}
}
