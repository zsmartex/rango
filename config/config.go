package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/cockroachdb/errors"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		NewConfig,
		provideConfig,
	),
)

type Rango struct {
	RbacSystem []string `env:"RANGO_RBAC_SYSTEM" envDefault:"admin,superadmin,operator"`
	RbacAdmin  []string `env:"RANGO_RBAC_ADMIN" envDefault:"admin,superadmin"`
}

type Config struct {
	HTTP            config.HTTP
	Kafka           config.Kafka
	Rango           Rango
	ApplicationName string `env:"APP_NAME" envDefault:"Rango"`
	JWTPublicKey    string `env:"JWT_PUBLIC_KEY"`
}

type ConfigOut struct {
	fx.Out

	HTTP            config.HTTP `name:"http_server"`
	Kafka           config.Kafka
	ApplicationName string `name:"application_name"`
}

func provideConfig(config *Config) ConfigOut {
	return ConfigOut{
		HTTP:            config.HTTP,
		Kafka:           config.Kafka,
		ApplicationName: config.ApplicationName,
	}
}

func NewConfig() (*Config, error) {
	conf := new(Config)

	if err := env.Parse(conf); err != nil {
		return nil, errors.Newf("parse config: %v", err)
	}

	return conf, nil
}
