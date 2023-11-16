package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
)

var Module = fx.Module("metrics",
	fx.Provide(
		NewMetrics,
	),
)

type Metrics struct {
	clients prometheus.Gauge
	subs    *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		clients: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "rango_hub_clients_count",
				Help: "Number of clients currently connected",
			},
		),
		subs: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rango_hub_subscriptions_count",
				Help: "Number of user subscribed to a topic",
			},
			[]string{"type", "topic"},
		),
	}
}

func (m *Metrics) RecordHubClientNew() {
	m.clients.Inc()
}

func (m *Metrics) RecordHubClientClose() {
	m.clients.Dec()
}

func (m *Metrics) RecordHubSubscription(typ, topic string) {
	m.subs.WithLabelValues(typ, topic).Inc()
}

func (m *Metrics) RecordHubUnsubscription(typ, topic string) {
	m.subs.WithLabelValues(typ, topic).Dec()
}
