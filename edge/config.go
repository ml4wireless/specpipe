package edge

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ml4wireless/specpipe/common"
	"github.com/spf13/viper"
)

type Config struct {
	Fm *FmConfig `mapstructure:"fm"`
	Iq *IqConfig `mapstructure:"iq"`
}

type FmConfig struct {
	Device  *DeviceConfig   `mapstructure:"device"`
	Rtlsdr  *RtlsdrFmConfig `mapstructure:"rtlsdr"`
	Nats    *NatsConfig     `mapstructure:"nats"`
	Logging *LoggingConfig  `mapstructure:"logging"`
}

type IqConfig struct {
	Device  *DeviceConfig   `mapstructure:"device"`
	Rtlsdr  *RtlsdrIqConfig `mapstructure:"rtlsdr"`
	Nats    *NatsConfig     `mapstructure:"nats"`
	Logging *LoggingConfig  `mapstructure:"logging"`
}

type DeviceConfig struct {
	Name      string
	Latitude  float32
	Longitude float32
}

type RtlsdrFmConfig struct {
	Freq         string
	SampleRate   string
	ResampleRate string
	Rpc          *RtlsdrRpcConfig `mapstructure:"rpc"`
}

type RtlsdrIqConfig struct {
	Freq       string
	SampleRate string
	Rpc        *RtlsdrRpcConfig `mapstructure:"rpc"`
}

type RtlsdrRpcConfig struct {
	ServerAddr string
	ServerPort string
}

type NatsConfig struct {
	Url     string
	Subject string
}

type LoggingConfig struct {
	Level string
}

func setFmConfigDefault() {
	viper.SetDefault("fm.device.name", watermill.NewShortUUID())
	viper.SetDefault("fm.device.latitude", 0.0)
	viper.SetDefault("fm.device.longitude", 0.0)

	viper.SetDefault("fm.rtlsdr.freq", "")
	viper.SetDefault("fm.rtlsdr.sampleRate", "170k")
	viper.SetDefault("fm.rtlsdr.resampleRate", "32k")
	viper.SetDefault("fm.rtlsdr.rpc.serverAddr", "127.0.0.1")
	viper.SetDefault("fm.rtlsdr.rpc.serverPort", "40000")

	viper.SetDefault("fm.nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("fm.nats.subject", "")

	viper.SetDefault("fm.logging.level", "info")
}

func setIqConfigDefault() {
	viper.SetDefault("iq.device.name", watermill.NewShortUUID())
	viper.SetDefault("iq.device.latitude", 0.0)
	viper.SetDefault("iq.device.longitude", 0.0)

	viper.SetDefault("iq.rtlsdr.freq", "")
	viper.SetDefault("iq.rtlsdr.sampleRate", "2048000")
	viper.SetDefault("iq.rtlsdr.rpc.serverAddr", "127.0.0.1")
	viper.SetDefault("iq.rtlsdr.rpc.serverPort", "40000")

	viper.SetDefault("iq.nats.url", "nats://127.0.0.1:4222")
	viper.SetDefault("iq.nats.subject", "")

	viper.SetDefault("iq.logging.level", "info")
}

func NewFmConfig() (*FmConfig, error) {
	setFmConfigDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	if c.Fm.Nats.Subject == "" {
		c.Fm.Nats.Subject = common.DataSubject(common.FM, c.Fm.Device.Name)
	}
	return c.Fm, nil
}

func NewIqConfig() (*IqConfig, error) {
	setIqConfigDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	if c.Iq.Nats.Subject == "" {
		c.Iq.Nats.Subject = common.DataSubject(common.IQ, c.Iq.Device.Name)
	}
	return c.Iq, nil
}
