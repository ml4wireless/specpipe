package server

import (
	"github.com/spf13/viper"
)

type Config struct {
	Controller *ControllerConfig `mapstructure:"controller"`
	Health     *HealthConfig     `mapstructure:"health"`
}

type ControllerConfig struct {
	Http    *HttpConfig    `mapstructure:"http"`
	Nats    *NatsConfig    `mapstructure:"nats"`
	Logging *LoggingConfig `mapstructure:"logging"`
}

type HttpConfig struct {
	ServerPort string
}

type NatsConfig struct {
	Url string
}

type HealthConfig struct {
	TimeoutSec int64
	Nats       *NatsConfig    `mapstructure:"nats"`
	Logging    *LoggingConfig `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level string
}

func setControllerConfigDefault() {
	viper.SetDefault("controller.http.serverPort", "8888")
	viper.SetDefault("controller.nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("controller.logging.level", "info")
}

func setHealthConfigDefault() {
	viper.SetDefault("health.timeoutSec", 5)
	viper.SetDefault("health.nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("health.logging.level", "info")
}

func NewControllerConfig() (*ControllerConfig, error) {
	setControllerConfigDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c.Controller, nil
}

func NewHealthConfig() (*HealthConfig, error) {
	setHealthConfigDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c.Health, nil
}
