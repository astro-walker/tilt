package engine

import (
	"github.com/astro-walker/tilt/internal/cloud"
	"github.com/astro-walker/tilt/internal/controllers"
	"github.com/astro-walker/tilt/internal/engine/analytics"
	"github.com/astro-walker/tilt/internal/engine/configs"
	"github.com/astro-walker/tilt/internal/engine/dockerprune"
	"github.com/astro-walker/tilt/internal/engine/k8srollout"
	"github.com/astro-walker/tilt/internal/engine/k8swatch"
	"github.com/astro-walker/tilt/internal/engine/local"
	"github.com/astro-walker/tilt/internal/engine/runtimelog"
	"github.com/astro-walker/tilt/internal/engine/session"
	"github.com/astro-walker/tilt/internal/engine/telemetry"
	"github.com/astro-walker/tilt/internal/engine/uiresource"
	"github.com/astro-walker/tilt/internal/engine/uisession"
	"github.com/astro-walker/tilt/internal/hud"
	"github.com/astro-walker/tilt/internal/hud/prompt"
	"github.com/astro-walker/tilt/internal/hud/server"
	"github.com/astro-walker/tilt/internal/store"
)

// Subscribers that only read from the new Tilt API,
// and run the API server.
func ProvideSubscribersAPIOnly(
	hudsc *server.HeadsUpServerController,
	tscm *controllers.TiltServerControllerManager,
	cb *controllers.ControllerBuilder,
	ts *hud.TerminalStream,
) []store.Subscriber {
	return []store.Subscriber{
		// The API server must go before other subscribers,
		// so that it can run its boot sequence first.
		hudsc,

		// The controller manager must go after the API server,
		// so that it can connect to it and make resources available.
		tscm,
		cb,
		ts,
	}
}

func ProvideSubscribers(
	hudsc *server.HeadsUpServerController,
	tscm *controllers.TiltServerControllerManager,
	cb *controllers.ControllerBuilder,
	hud hud.HeadsUpDisplay,
	ts *hud.TerminalStream,
	tp *prompt.TerminalPrompt,
	sw *k8swatch.ServiceWatcher,
	bc *BuildController,
	cc *configs.ConfigsController,
	tqs *configs.TriggerQueueSubscriber,
	dclm *runtimelog.DockerComposeLogManager,
	ar *analytics.AnalyticsReporter,
	au *analytics.AnalyticsUpdater,
	ewm *k8swatch.EventWatchManager,
	tcum *cloud.CloudStatusManager,
	dp *dockerprune.DockerPruner,
	tc *telemetry.Controller,
	lsc *local.ServerController,
	podm *k8srollout.PodMonitor,
	sc *session.Controller,
	uss *uisession.Subscriber,
	urs *uiresource.Subscriber,
) []store.Subscriber {
	apiSubscribers := ProvideSubscribersAPIOnly(hudsc, tscm, cb, ts)

	legacySubscribers := []store.Subscriber{
		hud,
		tp,
		sw,
		bc,
		cc,
		tqs,
		dclm,
		ar,
		au,
		ewm,
		tcum,
		dp,
		tc,
		lsc,
		podm,
		sc,
		uss,
		urs,
	}
	return append(apiSubscribers, legacySubscribers...)
}
