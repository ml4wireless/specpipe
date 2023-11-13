package server

import (
	"github.com/spf13/viper"
)

type Config struct {
	Http      *HttpConfig      `mapstructure:"http"`
	Nats      *NatsConfig      `mapstructure:"nats"`
	Heartbeat *HeartbeatConfig `mapstructure:"heartbeat"`
	Logging   *LoggingConfig   `mapstructure:"logging"`
}

type HttpConfig struct {
	ServerPort string
}

type NatsConfig struct {
	Url string
}

type HeartbeatConfig struct {
	TimeoutSec int64
}

type LoggingConfig struct {
	Level string
}

func setDefault() {
	viper.SetDefault("http.serverPort", "8888")
	viper.SetDefault("nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("heartbeat.timeoutSec", 5)
	viper.SetDefault("logging.level", "info")
}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
