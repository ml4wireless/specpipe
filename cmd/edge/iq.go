package main

import (
	"context"
	"errors"
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

var iqCmd = &cobra.Command{
	Use:   "iq",
	Short: "IQ data capturer",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := edge.NewIqConfig()
		if err != nil {
			log.Fatal(err)
		}

		publisher, err := common.NewNatsPublisher(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer publisher.Close()
		logger, err := common.NewEdgeLogrus(config.Logging.Level)
		if err != nil {
			log.Fatal(err)
		}
		clusterConn, err := common.NewNATSConn(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer clusterConn.Close()
		js, err := jetstream.New(clusterConn)
		if err != nil {
			log.Fatal(err)
		}
		regCtx, regCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer regCancel()
		kv, err := js.KeyValue(regCtx, common.KVStoreBucket)
		if err != nil {
			log.Fatal(err)
		}
		deviceInfo := common.IQDevice{
			Name:       config.Device.Name,
			Freq:       config.Rtlsdr.Freq,
			SampleRate: config.Rtlsdr.SampleRate,
			Latitude:   config.Device.Latitude,
			Longitude:  config.Device.Longitude,
		}
		heartbeatSub, err := edge.RegisterDevice(regCtx, clusterConn, kv, common.IQ, config.Device.Name, &deviceInfo)
		if err != nil {
			log.Fatal(err)
		}
		defer heartbeatSub.Unsubscribe()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		watchSub, iqDeviceChan, err := edge.WatchDevice(ctx, clusterConn, kv, common.IQ, config.Device.Name, logger)
		if err != nil {
			log.Fatal(err)
		}
		defer watchSub.Unsubscribe()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		done := make(chan bool, 1)
		go func() {
			subCtx, subCancel := context.WithCancel(ctx)

			for {
				select {
				case <-ctx.Done():
					subCancel()
					deregCtx, deregCancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer deregCancel()
					if err = edge.DeregisterDevice(deregCtx, kv, common.IQ, config.Device.Name); err != nil {
						log.Fatal(err)
					}
					logger.Info("device deregistered")
					done <- true
					return
				case newDevice := <-iqDeviceChan:
					newIqDevice, ok := newDevice.(*common.IQDevice)
					if !ok {
						logger.Error(errors.New("casting iq device type error"))
						break
					}
					subCancel()

					config.Rtlsdr.Freq = newIqDevice.Freq
					config.Rtlsdr.SampleRate = newIqDevice.SampleRate

					logger.Infof("device %s tuned to frequency=%s sampling_rate=%s", config.Device.Name, config.Rtlsdr.Freq, config.Rtlsdr.SampleRate)
					subCtx, subCancel = context.WithCancel(ctx)
					go func() {
						if err := edge.CaptureIQ(subCtx, config, publisher, logger); err != nil {
							log.Fatal(err)
						}
					}()
				}
			}
		}()

		log.Infof("start specpipe edge IQ name=%s nats-subject=%s", config.Device.Name, config.Nats.Subject)
		iqDeviceChan <- &deviceInfo

		<-sig
		cancel()

		<-done
	},
}

func init() {
	rootCmd.AddCommand(iqCmd)

	iqCmd.Flags().StringP("device-name", "", "", "device name")
	viper.BindPFlag("iq.device.name", iqCmd.Flags().Lookup("device-name"))

	iqCmd.Flags().StringP("device-latitude", "", "", "device latitude")
	viper.BindPFlag("iq.device.latitude", iqCmd.Flags().Lookup("device-latitude"))

	iqCmd.Flags().StringP("device-longitude", "", "", "device longitude")
	viper.BindPFlag("iq.device.longitude", iqCmd.Flags().Lookup("device-longitude"))

	iqCmd.Flags().StringP("freq", "", "", "frequency")
	viper.BindPFlag("iq.rtlsdr.freq", iqCmd.Flags().Lookup("freq"))

	iqCmd.Flags().StringP("sample-rate", "", "", "sampling rate")
	viper.BindPFlag("iq.rtlsdr.sampleRate", iqCmd.Flags().Lookup("sample-rate"))

	iqCmd.Flags().StringP("rpc-server-addr", "", "", "rtlsdr rpc server address")
	viper.BindPFlag("iq.rtlsdr.rpc.serverAddr", iqCmd.Flags().Lookup("rpc-server-addr"))

	iqCmd.Flags().StringP("rpc-server-port", "", "", "rtlsdr rpc server port")
	viper.BindPFlag("iq.rtlsdr.rpc.serverPort", iqCmd.Flags().Lookup("rpc-server-port"))

	iqCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("iq.nats.url", iqCmd.Flags().Lookup("nats-url"))

	iqCmd.Flags().StringP("nats-subject", "", "", "NATS subject")
	viper.BindPFlag("iq.nats.subject", iqCmd.Flags().Lookup("nats-subject"))

	iqCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("iq.logging.level", iqCmd.Flags().Lookup("log-level"))
}
