package fsevent

import (
	"time"

	"github.com/astro-walker/tilt/internal/watch"
	"github.com/astro-walker/tilt/pkg/logger"
)

type WatcherMaker func(paths []string, ignore watch.PathMatcher, l logger.Logger) (watch.Notify, error)

type TimerMaker func(d time.Duration) <-chan time.Time

func ProvideWatcherMaker() WatcherMaker {
	return watch.NewWatcher
}

func ProvideTimerMaker() TimerMaker {
	return time.After
}
