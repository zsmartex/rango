package routing

import "go.uber.org/fx"

var Module = fx.Module("routing",
	fx.Provide(
		NewHub,
		NewTopic,
	),
)
