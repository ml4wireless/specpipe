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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var forwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "IQ data forwarder",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := edge.NewForwardConfig()
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

		watchSub, forwardDeviceChan, err := edge.WatchDevice(ctx, clusterConn, kv, common.IQ, config.Device.Name, logger)
		if err != nil {
			log.Fatal(err)
		}
		defer watchSub.Unsubscribe()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		conn, err := initGrpcClient(ctx, config.Grpc.Target)

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
				case newDevice := <-forwardDeviceChan:
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
						if err := edge.ForwardIQ(subCtx, config, logger, conn); err != nil {
							log.Fatal(err)
						}
					}()
				}
			}
		}()

		log.Infof("start specpipe edge IQ name=%s nats-subject=%s", config.Device.Name, config.Nats.Subject)
		forwardDeviceChan <- &deviceInfo

		<-sig
		cancel()

		<-done
	},
}

func initGrpcClient(ctx context.Context, grpcTarget string) (*grpc.ClientConn, error) {
	tlsTransCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, grpcTarget, tlsTransCreds)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func init() {
	rootCmd.AddCommand(forwardCmd)

	forwardCmd.Flags().StringP("target", "", "", "grpc server target")
	viper.BindPFlag("forward.grpc.target", forwardCmd.Flags().Lookup("target"))

	forwardCmd.Flags().StringP("device-name", "", "", "device name")
	viper.BindPFlag("forward.device.name", forwardCmd.Flags().Lookup("device-name"))

	forwardCmd.Flags().StringP("device-latitude", "", "", "device latitude")
	viper.BindPFlag("forward.device.latitude", forwardCmd.Flags().Lookup("device-latitude"))

	forwardCmd.Flags().StringP("device-longitude", "", "", "device longitude")
	viper.BindPFlag("forward.device.longitude", forwardCmd.Flags().Lookup("device-longitude"))

	forwardCmd.Flags().StringP("freq", "", "", "frequency")
	viper.BindPFlag("forward.rtlsdr.freq", forwardCmd.Flags().Lookup("freq"))

	forwardCmd.Flags().StringP("sample-rate", "", "", "sampling rate")
	viper.BindPFlag("forward.rtlsdr.sampleRate", forwardCmd.Flags().Lookup("sample-rate"))

	forwardCmd.Flags().StringP("rpc-server-addr", "", "", "rtlsdr rpc server address")
	viper.BindPFlag("forward.rtlsdr.rpc.serverAddr", forwardCmd.Flags().Lookup("rpc-server-addr"))

	forwardCmd.Flags().StringP("rpc-server-port", "", "", "rtlsdr rpc server port")
	viper.BindPFlag("forward.rtlsdr.rpc.serverPort", forwardCmd.Flags().Lookup("rpc-server-port"))

	forwardCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("forward.nats.url", forwardCmd.Flags().Lookup("nats-url"))

	forwardCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("forward.logging.level", forwardCmd.Flags().Lookup("log-level"))
}
