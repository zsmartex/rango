package app

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/zsmartex/pkg/v2/infrastructure/kafka_fx"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"

	"github.com/zsmartex/rango/config"
	"github.com/zsmartex/rango/pkg/routing"
)

var Topic kafka_fx.Topic = "rango.events"
var Ticker *time.Ticker = time.NewTicker(20 * time.Millisecond)

var (
	Module = fx.Module("rango_service",
		kafka_fx.ConsumerModule,
		rangoProviders,
	)

	rangoProviders = fx.Options(
		fx.Supply(Topic),
		fx.Supply(Ticker),
		fx.Supply(fx.Annotate(true, fx.ParamTags(`name:"at_end"`))),
		fx.Provide(
			New,
			func(app *App) kafka_fx.ConsumerSubscriber {
				return app
			},
		),
		fx.Invoke(registerHooks),
	)
)

var _ kafka_fx.ConsumerSubscriber = (*App)(nil)

type App struct {
	config *config.Config
	hub    *routing.Hub
}

func New(config *config.Config, hub *routing.Hub) *App {
	app := &App{
		config: config,
		hub:    hub,
	}

	return app
}

func (s *App) OnMessage(record *kgo.Record) error {
	s.hub.ReceiveMsg(record)

	return nil
}

func (a *App) Run() error {
	pub, err := getPublicKey(a.config)
	if err != nil {
		time.Sleep(2 * time.Second)
		return errors.Wrap(err, "Loading public key failed")
	}

	log.Info("Starting rango...")

	go a.hub.ListenWebsocketEvents()

	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		routing.NewClient(a.hub, w, r)
	}

	http.HandleFunc("/private", authHandler(wsHandler, pub, true))
	http.HandleFunc("/public", authHandler(wsHandler, pub, false))
	http.HandleFunc("/", authHandler(wsHandler, pub, false))

	go func() {
		err := http.ListenAndServe(":4242", promhttp.Handler())
		if err != nil {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	log.Infof("Listenning on %s", a.config.HTTP.Address())
	err = http.ListenAndServe(a.config.HTTP.Address(), nil)
	if err != nil {
		log.Errorf("ListenAndServe failed: %v", err)
		return err
	}

	return nil
}

func registerHooks(app *App) {
	go app.Run()
}
