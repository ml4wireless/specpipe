package edge

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/spf13/viper"
)

type Config struct {
	Device  *DeviceConfig  `mapstructure:"device"`
	Rtlsdr  *RtlsdrConfig  `mapstructure:"rtlsdr"`
	Pub     *PublishConfig `mapstructure:"pub"`
	Logging *LoggingConfig `mapstructure:"logging"`
}

type DeviceConfig struct {
	Name      string
	Longitude float32
	Latitude  float32
}

type RtlsdrConfig struct {
	Fm struct {
		Freq         string
		SampleRate   string
		ResampleRate string
	}
	RpcServerAddr string
	RpcServerPort string
}

type PublishConfig struct {
	NatsUrl string
	Subject string
}

type LoggingConfig struct {
	Level string
}

func setDefault() {
	viper.SetDefault("device.name", watermill.NewShortUUID())
	viper.SetDefault("device.longitude", 0.0)
	viper.SetDefault("device.latitude", 0.0)
	viper.SetDefault("rtlsdr.fm.freq", "")
	viper.SetDefault("rtlsdr.fm.sampleRate", "170k")
	viper.SetDefault("rtlsdr.fm.resampleRate", "32k")
	viper.SetDefault("rtlsdr.rpcServerAddr", "127.0.0.1")
	viper.SetDefault("rtlsdr.rpcServerPort", "40000")
	viper.SetDefault("pub.natsUrl", "nats://127.0.0.1:4222")
	viper.SetDefault("pub.subject", "specpipe.fm")
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
