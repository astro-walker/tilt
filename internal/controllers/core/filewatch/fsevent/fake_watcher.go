package fsevent

import (
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/astro-walker/tilt/internal/ospath"
	"github.com/astro-walker/tilt/internal/watch"
	"github.com/astro-walker/tilt/pkg/logger"
)

type FakeMultiWatcher struct {
	Events chan watch.FileEvent
	Errors chan error

	mu         sync.Mutex
	watchers   []*FakeWatcher
	subs       []chan watch.FileEvent
	subsErrors []chan error
}

func NewFakeMultiWatcher() *FakeMultiWatcher {
	r := &FakeMultiWatcher{
		Events: make(chan watch.FileEvent, 20),
		Errors: make(chan error, 20),
	}
	go r.loop()
	return r
}

func (w *FakeMultiWatcher) NewSub(paths []string, ignore watch.PathMatcher, _ logger.Logger) (watch.Notify, error) {
	subCh := make(chan watch.FileEvent)
	errorCh := make(chan error)
	w.mu.Lock()
	defer w.mu.Unlock()

	watcher := NewFakeWatcher(subCh, errorCh, paths, ignore)
	w.watchers = append(w.watchers, watcher)
	w.subs = append(w.subs, subCh)
	w.subsErrors = append(w.subsErrors, errorCh)
	return watcher, nil
}

func (w *FakeMultiWatcher) getSubs() []chan watch.FileEvent {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]chan watch.FileEvent{}, w.subs...)
}

func (w *FakeMultiWatcher) getSubErrors() []chan error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]chan error{}, w.subsErrors...)
}

func (w *FakeMultiWatcher) loop() {
	defer func() {
		for _, sub := range w.getSubs() {
			close(sub)
		}

		for _, sub := range w.getSubErrors() {
			close(sub)
		}
	}()

	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				return
			}
			w.mu.Lock()
			for _, watcher := range w.watchers {
				if watcher.Running && watcher.matches(e.Path()) {
					watcher.inboundCh <- e
				}
			}
			w.mu.Unlock()
		case e, ok := <-w.Errors:
			if !ok {
				return
			}
			for _, sub := range w.getSubErrors() {
				sub <- e
			}
		}
	}
}

type FakeWatcher struct {
	inboundCh  chan watch.FileEvent
	outboundCh chan watch.FileEvent
	errorCh    chan error
	closeCh    chan bool

	eventCount uint64

	paths  []string
	ignore watch.PathMatcher

	Running  bool
	StartErr error
}

func NewFakeWatcher(inboundCh chan watch.FileEvent, errorCh chan error, paths []string, ignore watch.PathMatcher) *FakeWatcher {
	for i, path := range paths {
		paths[i], _ = filepath.Abs(path)
	}

	return &FakeWatcher{
		inboundCh:  inboundCh,
		outboundCh: make(chan watch.FileEvent, 20),
		errorCh:    errorCh,
		paths:      paths,
		ignore:     ignore,
		closeCh:    make(chan bool),
	}
}

func (w *FakeWatcher) matches(path string) bool {
	ignore, _ := w.ignore.Matches(path)
	if ignore {
		return false
	}

	for _, watched := range w.paths {
		if ospath.IsChild(watched, path) {
			return true
		}
	}
	return false
}

func (w *FakeWatcher) Start() error {
	w.Running = true
	go w.loop()
	if w.StartErr != nil {
		return w.StartErr
	}
	return nil
}

func (w *FakeWatcher) Close() error {
	close(w.closeCh)
	<-w.outboundCh
	return nil
}

func (w *FakeWatcher) Errors() chan error {
	return w.errorCh
}

func (w *FakeWatcher) Events() chan watch.FileEvent {
	return w.outboundCh
}

func (w *FakeWatcher) TotalEventCount() uint64 {
	return atomic.LoadUint64(&w.eventCount)
}

func (w *FakeWatcher) QueuedCount() int {
	return len(w.outboundCh)
}

func (w *FakeWatcher) loop() {
	defer func() {
		w.Running = false
		close(w.outboundCh)
	}()

	var q []watch.FileEvent
	for {
		if len(q) == 0 {
			select {
			case e, ok := <-w.inboundCh:
				if !ok {
					return
				}
				q = append(q, e)
			case <-w.closeCh:
				return
			}
		} else {
			e := q[0]
			w.outboundCh <- e
			atomic.AddUint64(&w.eventCount, 1)
			q = q[1:]
		}
	}
}

var _ watch.Notify = &FakeWatcher{}
