package config

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = ""

type Config struct {
	Prometheus Prometheus
	Grafana    Grafana
	Tempo      Tempo
}

type Prometheus struct {
	TmplName string `envconfig:"PROMETHEUS_TEMPLATE_NAME" default:"Prometheus"`
	Version  string `envconfig:"PROMETHEUS_VERSION" default:"latest"`
}

type Grafana struct {
	TmplName string `envconfig:"GRAFANA_TEMPLATE_NAME" default:"Grafana"`
	Version  string `envconfig:"GRAFANA_VERSION" default:"latest"`
}

type Tempo struct {
	TmplName string `envconfig:"TEMPO_TEMPLATE_NAME" default:"tempo"`
	Version  string `envconfig:"TEMPO_VERSION" default:"latest"`
}

func LoadConf() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
