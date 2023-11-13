package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ml4wireless/specpipe/common"
	"github.com/ml4wireless/specpipe/edge"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fmCmd = &cobra.Command{
	Use:   "fm",
	Short: "FM audio capturer",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := edge.NewConfig()
		if err != nil {
			log.Fatal(err)
		}

		publisher, err := common.NewNatsPublisher(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		logger, err := common.NewEdgeLogrus(config.Logging.Level)
		if err != nil {
			log.Fatal(err)
		}
		clusterConn, err := common.NewNATSConn(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer clusterConn.Close()
		kvConn, err := common.NewNATSConn(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer kvConn.Close()
		js, err := jetstream.New(kvConn)
		if err != nil {
			log.Fatal(err)
		}
		regCtx, regCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer regCancel()
		kv, err := js.KeyValue(regCtx, common.KVStoreBucket)
		if err != nil {
			log.Fatal(err)
		}
		deviceInfo := common.FMDevice{
			Name:       config.Device.Name,
			Freq:       config.Rtlsdr.Fm.Freq,
			SampleRate: config.Rtlsdr.Fm.SampleRate,
			Latitude:   config.Device.Latitude,
			Longitude:  config.Device.Longitude,
		}
		deviceInfoBytes, err := json.Marshal(deviceInfo)
		if err != nil {
			log.Fatal(err)
		}
		heartbeatSub, err := edge.RegisterDevice(regCtx, clusterConn, kv, common.FM, config.Device.Name, deviceInfoBytes)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		log.Infof("start specpipe edge FM name=%s freq=%s nats-subject=%s", config.Device.Name, config.Rtlsdr.Fm.Freq, config.Nats.Subject)
		done := make(chan bool, 1)
		go func() {
			if err := edge.CaptureAudio(ctx, config, publisher, logger); err != nil {
				log.Fatal(err)
			}
			logger.Info("audio capture stopped")
			deregCtx, deregCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer deregCancel()
			if err = edge.DeregisterDevice(deregCtx, kv, common.FM, config.Device.Name, heartbeatSub); err != nil {
				log.Fatal(err)
			}
			logger.Info("device deregistered")
			done <- true
		}()
		<-sig
		cancel()

		<-done
	},
}

func init() {
	rootCmd.AddCommand(fmCmd)

	fmCmd.Flags().StringP("device-name", "", "", "device name")
	viper.BindPFlag("device.name", fmCmd.Flags().Lookup("device-name"))

	fmCmd.Flags().StringP("device-latitude", "", "", "device latitude")
	viper.BindPFlag("device.latitude", fmCmd.Flags().Lookup("device-latitude"))

	fmCmd.Flags().StringP("device-longitude", "", "", "device longitude")
	viper.BindPFlag("device.longitude", fmCmd.Flags().Lookup("device-longitude"))

	fmCmd.Flags().StringP("freq", "", "", "frequency")
	viper.BindPFlag("rtlsdr.fm.freq", fmCmd.Flags().Lookup("freq"))

	fmCmd.Flags().StringP("sample-rate", "", "", "sampling rate")
	viper.BindPFlag("rtlsdr.fm.sampleRate", fmCmd.Flags().Lookup("sample-rate"))

	fmCmd.Flags().StringP("resample-rate", "", "", "sampling rate")
	viper.BindPFlag("rtlsdr.fm.resampleRate", fmCmd.Flags().Lookup("resample-rate"))

	fmCmd.Flags().StringP("rpc-server-addr", "", "", "rtlsdr rpc server address")
	viper.BindPFlag("rtlsdr.rpcServerAddr", fmCmd.Flags().Lookup("rpc-server-addr"))

	fmCmd.Flags().StringP("rpc-server-port", "", "", "rtlsdr rpc server port")
	viper.BindPFlag("rtlsdr.rpcServerPort", fmCmd.Flags().Lookup("rpc-server-port"))

	fmCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("nats.url", fmCmd.Flags().Lookup("nats-url"))

	fmCmd.Flags().StringP("nats-subject", "", "", "NATS subject")
	viper.BindPFlag("nats.subject", fmCmd.Flags().Lookup("nats-subject"))

	fmCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("logging.level", fmCmd.Flags().Lookup("log-level"))
}
