package main

import (
	"go.uber.org/fx"

	"github.com/zsmartex/rango/cmd/rango/app.go"
	"github.com/zsmartex/rango/config"
	"github.com/zsmartex/rango/pkg/metrics"
	"github.com/zsmartex/rango/pkg/routing"
)

func main() {
	app := fx.New(
		config.Module,
		metrics.Module,
		routing.Module,
		app.Module,
	)

	app.Run()
}
