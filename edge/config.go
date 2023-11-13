package edge

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ml4wireless/specpipe/common"
	"github.com/spf13/viper"
)

type Config struct {
	Device  *DeviceConfig  `mapstructure:"device"`
	Rtlsdr  *RtlsdrConfig  `mapstructure:"rtlsdr"`
	Nats    *NatsConfig    `mapstructure:"nats"`
	Logging *LoggingConfig `mapstructure:"logging"`
}

type DeviceConfig struct {
	Name      string
	Latitude  float32
	Longitude float32
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

type NatsConfig struct {
	Url     string
	Subject string
}

type LoggingConfig struct {
	Level string
}

func setDefault() {
	viper.SetDefault("device.name", watermill.NewShortUUID())
	viper.SetDefault("device.latitude", 0.0)
	viper.SetDefault("device.longitude", 0.0)
	viper.SetDefault("rtlsdr.fm.freq", "")
	viper.SetDefault("rtlsdr.fm.sampleRate", "170k")
	viper.SetDefault("rtlsdr.fm.resampleRate", "32k")
	viper.SetDefault("rtlsdr.rpcServerAddr", "127.0.0.1")
	viper.SetDefault("rtlsdr.rpcServerPort", "40000")
	viper.SetDefault("nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("logging.level", "info")
}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	c.Nats.Subject = common.DataSubject(common.FM, c.Device.Name)
	return &c, nil
}
